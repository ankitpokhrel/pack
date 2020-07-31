package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/git-lfs/wildmatch"
	"github.com/mholt/archiver/v3"
	"github.com/urfave/cli/v2"
)

type list map[string]struct{}

var ignoreList cli.StringSlice

func main() {
	app := &cli.App{
		Name:  "pack",
		Usage: "Pack compresses file/folder ignoring any hidden and given files",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "ignore",
				Aliases:     []string{"ig"},
				Usage:       "Ignore list from given files",
				Destination: &ignoreList,
			},
		},
		Action: pack,
	}

	if err := app.Run(os.Args); err != nil && err != io.EOF {
		panic(err)
	}
}

func pack(c *cli.Context) error {
	src, dest := ".", defaultFileName()
	if c.NArg() > 0 {
		src, dest = c.Args().Get(0), c.Args().Get(1)
	}

	accept, err := files(src, ignoreList.Value())
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempDir("", "pack")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	rootDir := accept[0]
	basePath, err := createRootDir(dest, tmp)
	if err != nil {
		return err
	}

	for _, v := range accept {
		if v == rootDir {
			continue
		}

		newPath := basePath + strings.Replace(v, rootDir, "", 1)

		if err := dupe(v, newPath); err != nil {
			return err
		}
	}

	return archiver.Archive([]string{basePath}, dest)
}

func files(path string, avoid []string) ([]string, error) {
	var accept []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = filepath.Clean(path)

		if isHidden(path) {
			if isDir(path) {
				return filepath.SkipDir
			}
			return nil
		}

		for _, pattern := range avoid {
			wm := wildmatch.NewWildmatch(pattern, wildmatch.Basename)
			if wm.Match(path) {
				if isDir(path) {
					return filepath.SkipDir
				}
				return nil
			}
		}

		accept = append(accept, path)

		return nil
	})

	return accept, err
}

func ignore(files []string) list {
	var avoid = make(list)

	for _, file := range files {
		path := parsePath(file)

		func() {
			f, err := os.Open(path)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				avoid[scanner.Text()] = struct{}{}
			}
		}()
	}

	return avoid
}

func isHidden(path string) bool {
	if len(path) == 0 {
		return true
	}

	base := filepath.Base(path)

	return len(base) > 0 && base[0] == '.'
}

func parsePath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Clean(strings.ReplaceAll(path, "~", homeDir))
}

func defaultFileName() string {
	var name strings.Builder

	name.WriteString("archive-")
	name.WriteString(strconv.Itoa((int)(time.Now().Unix())))
	name.WriteString(".zip")

	return name.String()
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func createRootDir(dest, path string) (string, error) {
	name, ext := filepath.Base(dest), filepath.Ext(dest)

	name = strings.Replace(name, ext, "", 1)

	basePath := path + "/" + name

	return basePath, exec.Command("mkdir", "-p", basePath).Run()
}

func dupe(src, dest string) error {
	if isDir(src) {
		return exec.Command("mkdir", "-p", dest).Run()
	}
	return exec.Command("cp", "-rf", src, dest).Run()
}
