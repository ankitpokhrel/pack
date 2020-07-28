package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	accept, err := files(src)
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempDir("", "pack")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	rootDir := accept[0]
	basePath, err := createMainDir(dest, tmp)
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

func files(path string) ([]string, error) {
	var accept []string

	avoid := ignore()
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

		for pattern := range avoid {
			matched, err := regexp.MatchString(pattern, path)
			if err != nil {
				return err
			}
			if matched {
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

func ignore() list {
	var avoid = make(list)

	for _, file := range ignoreList.Value() {
		path := parsePath(file)

		func() {
			f, err := os.Open(path)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				avoid[filepath.Clean(scanner.Text())] = struct{}{}
			}
		}()
	}

	return avoid
}

// Only for unix based systems.
func isHidden(path string) bool {
	if len(path) == 0 {
		return false
	}

	base := filepath.Base(path)

	return len(base) > 0 && base[0] == '.'
}

func parsePath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return strings.ReplaceAll(path, "~", homeDir)
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

func createMainDir(dest, path string) (string, error) {
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
