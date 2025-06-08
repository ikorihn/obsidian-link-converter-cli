package olconv

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func LinkToWikilink(basepath string) error {

	files, err := ListMdFiles(basepath)
	if err != nil {
		panic(err)
	}
	filemap := FileListToMap(files)
	c := NewConverter(filemap)

	for _, file := range files {
		newLineAtEnd := false
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if len(content) == 0 {
			continue
		}
		if content[len(content)-1] == '\n' {
			newLineAtEnd = true
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		if err := c.Convert(f, buf, newLineAtEnd, ToWikilink); err != nil {
			return err
		}
		f.Close()

		wf, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		wf.WriteString(buf.String())
		wf.Close()
	}

	return nil

}

func WikilinkToLink(basepath string) error {
	files, err := ListMdFiles(basepath)
	if err != nil {
		return err
	}
	filemap := FileListToMap(files)
	c := NewConverter(filemap)

	for _, file := range files {
		newLineAtEnd := false
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if len(content) == 0 {
			continue
		}
		if content[len(content)-1] == '\n' {
			newLineAtEnd = true
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		if err := c.Convert(f, buf, newLineAtEnd, ToMarkdown); err != nil {
			return err
		}
		f.Close()

		wf, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		wf.WriteString(buf.String())
		wf.Close()
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
	filename, _ = url.QueryUnescape(filename)

	return filename
}
func formatRelativePath(relativePath string) string {
	relativePath = strings.TrimPrefix(relativePath, "./")
	relativePath = strings.TrimSuffix(relativePath, ".md")
	relativePath, _ = url.QueryUnescape(relativePath)
	return relativePath
}
