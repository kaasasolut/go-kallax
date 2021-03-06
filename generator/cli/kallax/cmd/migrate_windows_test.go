// +build windows

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathToFileURL(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	cases := []struct {
		input    string
		expected string
	}{
		{"c:\\foo\\bar\\baz", "file://c:/foo/bar/baz"},
		{"foo\\bar", "file://" + filepath.ToSlash(filepath.Join(wd, "foo", "bar"))},
	}

	for _, tt := range cases {
		require.Equal(t, tt.expected, pathToFileURL(tt.input), tt.input)
	}
}
