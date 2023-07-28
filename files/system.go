package files

import (
	"fmt"
	"sort"

	"go.uber.org/zap/zapcore"
)

type Validatable interface {
	Validate() error
}

type System[T Validatable] struct {
	files map[Path]T
}

func NewSystem[T Validatable](in map[Path]T) *System[T] {
	return &System[T]{files: in}
}

func (f *System[T]) Add(path Path, state T) error {
	if f.files == nil {
		f.files = make(map[Path]T)
	}
	path = path.Clean()
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if err := state.Validate(); err != nil {
		return fmt.Errorf("invalid state for %s: %w", path, err)
	}
	if _, ok := f.files[path]; ok {
		return fmt.Errorf("file %s already exists", path)
	}
	f.files[path] = state
	return nil
}

func (f *System[T]) MarshalLogObject(e zapcore.ObjectEncoder) error {
	if err := e.AddArray("paths", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, path := range f.Paths() {
			file := f.Get(path)
			if err := enc.AppendObject(zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				enc.AddString("path", string(path))
				if err := enc.AddReflected("state", file); err != nil {
					return fmt.Errorf("failed to marshal state: %w", err)
				}
				return nil
			})); err != nil {
				return fmt.Errorf("failed to marshal path %s: %w", path, err)
			}
		}
		return nil
	})); err != nil {
		return fmt.Errorf("failed to marshal paths: %w", err)
	}
	return nil
}

var _ zapcore.ObjectMarshaler = &System[Validatable]{}

func (f *System[T]) Paths() []Path {
	paths := make([]Path, 0, len(f.files))
	for path := range f.files {
		paths = append(paths, path)
	}
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})
	return paths
}

func (f *System[T]) Get(path Path) T {
	path = path.Clean()
	return f.files[path]
}

func (f *System[T]) Remove(path Path) (T, bool) {
	path = path.Clean()
	if f.IsTracked(path) {
		ret := f.files[path]
		delete(f.files, path)
		return ret, true
	}
	var ret T
	return ret, false
}

func (f *System[T]) IsTracked(path Path) bool {
	path = path.Clean()
	if f == nil || f.files == nil {
		return false
	}
	_, ok := f.files[path]
	return ok
}

func (f *System[T]) RemoveTracked(path Path) error {
	path = path.Clean()
	if f.files == nil {
		return fmt.Errorf("file %s does not exist", path)
	}
	if _, ok := f.files[path]; !ok {
		return fmt.Errorf("file %s does not exist", path)
	}
	delete(f.files, path)
	return nil
}

func (f *System[T]) RemoveAll(paths []Path) {
	for _, path := range paths {
		path = path.Clean()
		if f.IsTracked(path) {
			delete(f.files, path)
		}
	}
}

type MergeDuplicatePathErr[T Validatable] struct {
	Path   Path
	Value1 T
	Value2 T
}

func (e *MergeDuplicatePathErr[T]) Error() string {
	return fmt.Sprintf("duplicate path %s: %v %v", e.Path, e.Value1, e.Value2)
}

func SystemMerge[T Validatable](systems ...*System[T]) (*System[T], error) {
	var ret System[T]
	for _, system := range systems {
		for path, state := range system.files {
			if ret.IsTracked(path) {
				return nil, &MergeDuplicatePathErr[T]{Path: path, Value1: ret.Get(path), Value2: state}
			}
			if err := ret.Add(path, state); err != nil {
				return nil, fmt.Errorf("failed to add %s: %w", path, err)
			}
		}
	}
	return &ret, nil
}

func ConvertToRemovals(s *System[*State]) (*System[*StateWithChangeReason], error) {
	var ret System[*StateWithChangeReason]
	for _, path := range s.Paths() {
		if err := ret.Add(path, &StateWithChangeReason{
			State: State{
				FileExistence: FileExistenceAbsent,
			},
			ChangeReason: &ChangeReason{
				Reason: "no-longer-tracked",
			},
		}); err != nil {
			return nil, fmt.Errorf("failed to add %s: %w", path, err)
		}
	}
	return &ret, nil
}
