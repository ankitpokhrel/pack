package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_files(t *testing.T) {
	t.Parallel()

	accept, err := files("./testdata/zipme", []string{"*.png", "*.pdf", "**/folder2/"})

	assert.Nil(t, err)
	assert.Equal(t, len(accept), 9)
}

func Test_ignore(t *testing.T) {
	t.Parallel()

	avoid := ignore([]string{"./testdata/zipme/.ignore", "./testdata/zipme/.ignore1", "./testdata/unknown"})
	keys := []string{"*.png", "*.pdf", "folder2/"}

	assert.Equal(t, 3, len(avoid))

	for _, k := range keys {
		_, ok := avoid[k]
		assert.True(t, ok)
	}
}

func Test_isHidden(t *testing.T) {
	t.Parallel()

	cases := []struct {
		file     string
		expected bool
	}{
		{".hidden", true},
		{"regular", false},
		{"nothidden.", false},
		{"nor.mal", false},
		{"..hidden", true},
		{".", true},
		{"..", true},
		{"", true},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, isHidden(tc.file))
	}
}

func Test_parsePath(t *testing.T) {
	t.Parallel()

	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	cases := []struct {
		file, expected string
	}{
		{"/some/path/to/file.ext", "/some/path/to/file.ext"},
		{"~/Desktop", homeDir + "/Desktop"},
		{"/some/path/./to/../file.ext", "/some/path/file.ext"},
		{"some/path/./with/dot.ext", "some/path/with/dot.ext"},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, parsePath(tc.file))
	}
}

func Test_defaultFileName(t *testing.T) {
	t.Parallel()

	assert.Regexp(t, "archive-[0-9]+.zip", defaultFileName())
}

func Test_isDir(t *testing.T) {
	t.Parallel()

	assert.True(t, isDir("./testdata"))
	assert.False(t, isDir("./testdata/empty.txt"))
}

func Test_createRootDir(t *testing.T) {
	t.Parallel()

	tmp := "/tmp/pack"

	path, err := createRootDir("./testdata/empty.txt", tmp)

	assert.Nil(t, err)
	assert.Equal(t, tmp+"/empty", path)
	assert.True(t, isDir(path))

	_ = os.RemoveAll("/tmp/pack")
}

func Test_dupe(t *testing.T) {
	t.Parallel()

	src, dest, file := "./testdata", "/tmp/pack-dupe", "/empty.txt"

	err := dupe(src, dest)
	assert.Nil(t, err)
	assert.True(t, isDir(dest))

	err = dupe(src+file, dest+file)
	assert.Nil(t, err)

	_, err = os.Stat(dest + file)
	assert.False(t, isDir(dest+file))
	assert.True(t, !os.IsNotExist(err))

	_ = os.RemoveAll("/tmp/pack-dupe")
}
