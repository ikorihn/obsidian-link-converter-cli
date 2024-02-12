package olconv

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func ConvertUnderDir(basepath string) error {

	files, err := ListMdFiles(basepath)
	if err != nil {
		panic(err)
	}
	filemap := FileListToMap(files)
	c := NewConverter(filemap)

	var eg errgroup.Group
	for _, file := range files {
		eg.Go(func() error {
			f, err := os.Open(file)
			if err != nil {
				panic(err)
			}

			buf := &bytes.Buffer{}
			if err := c.Convert(f, buf); err != nil {
				panic(err)
			}
			f.Close()

			wf, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				panic(err)
			}
			wf.WriteString(buf.String())

			wf.Close()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil

}

func ListMdFiles(basepath string) ([]string, error) {
	filelist := make([]string, 0)

	err := filepath.Walk(basepath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			switch f.Name() {
			case ".git", ".obsidian", ".trash":
				return filepath.SkipDir
			default:
				return nil
			}
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		filelist = append(filelist, path)

		return nil
	})

	return filelist, err
}

func FileListToMap(filelist []string) map[string][]string {
	filemap := make(map[string][]string)
	for _, path := range filelist {
		filename := filenameWithoutMdExtension(path)
		if list, ok := filemap[filename]; !ok {
			filemap[filename] = []string{path}
		} else {
			filemap[filename] = append(list, path)
		}
	}
	return filemap
}

func filenameWithoutMdExtension(fullpath string) string {
	filename := filepath.Base(fullpath)
	filename = strings.ReplaceAll(filename, ".md", "")
	return filename
}
func formatRelativePath(relativePath string) string {
	relativePath = strings.TrimPrefix(relativePath, "./")
	relativePath = strings.TrimSuffix(relativePath, ".md")
	return relativePath
}
