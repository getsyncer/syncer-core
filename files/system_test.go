package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystem(t *testing.T) {
	var f System[*State]
	require.Empty(t, f.Paths())
	require.False(t, f.IsTracked("foo"))
	require.Nil(t, f.Get("foo"))
	require.Error(t, f.Add("foo", &State{Contents: []byte("foo")}))
	s := State{FileExistence: FileExistencePresent, Mode: os.FileMode(0644), Contents: []byte("foo")}
	require.NoError(t, f.Add("foo", &s))
	require.Equal(t, []Path{"foo"}, f.Paths())
	require.True(t, f.IsTracked("foo"))
	require.Equal(t, &s, f.Get("foo"))
	require.NoError(t, f.RemoveTracked("foo"))
	require.Empty(t, f.Paths())
}

func TestSystemMerge(t *testing.T) {
	s1 := State{FileExistence: FileExistencePresent, Mode: os.FileMode(0644), Contents: []byte("foo")}
	s2 := State{FileExistence: FileExistencePresent, Mode: os.FileMode(0644), Contents: []byte("bar")}
	f1 := NewSystem[*State](map[Path]*State{
		"foo": &s1,
	})
	f2 := NewSystem[*State](map[Path]*State{
		"bar": &s2,
	})
	f3, err := SystemMerge[*State](f1, f2)
	require.NoError(t, err)
	require.Equal(t, []Path{"bar", "foo"}, f3.Paths())
	require.Equal(t, &s1, f3.Get("foo"))
	require.Equal(t, &s2, f3.Get("bar"))
}

func TestConvertToRemovals(t *testing.T) {
	s := State{FileExistence: FileExistencePresent, Mode: os.FileMode(0644), Contents: []byte("foo")}
	f := NewSystem[*State](map[Path]*State{
		"foo": &s,
	})
	f2, err := ConvertToRemovals(f)
	require.NoError(t, err)
	require.Equal(t, []Path{"foo"}, f2.Paths())
	require.Equal(t, &StateWithChangeReason{
		State: State{
			FileExistence: FileExistenceAbsent,
		},
		ChangeReason: &ChangeReason{
			Reason: "no-longer-tracked",
		},
	}, f2.Get("foo"))
}
