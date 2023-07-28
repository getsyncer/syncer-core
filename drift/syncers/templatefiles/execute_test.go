package templatefiles

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTemplate(t *testing.T) {
	tp, err := newTemplate("test", "Hello world")
	require.NoError(t, err)
	require.NotNil(t, tp)
	var into bytes.Buffer
	require.NoError(t, tp.Execute(&into, nil))
	require.Equal(t, "Hello world", into.String())
}
