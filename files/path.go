package files

import "path/filepath"

type Path string

func (f Path) Clean() Path {
	return Path(filepath.Clean(string(f)))
}

func (f Path) String() string {
	return string(f)
}
