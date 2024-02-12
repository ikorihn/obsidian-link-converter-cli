package olconv

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

var mdLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)

type Converter struct {
	inCodeBlock bool

	filemap map[string][]string
}

func NewConverter(filemap map[string][]string) *Converter {
	return &Converter{
		filemap: filemap,
	}
}

func (c *Converter) Convert(r io.Reader, w io.Writer) error {

	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)

	for sc.Scan() {
		line := sc.Text()

		line = c.convertLine(line)
		bw.WriteString(line)
		bw.WriteString("\n")

	}
	bw.Flush()

	return nil
}

func (c *Converter) convertLine(line string) string {
	if strings.HasPrefix(line, "```") {
		c.inCodeBlock = !c.inCodeBlock
	}
	if c.inCodeBlock {
		return line
	}

	matches := mdLinkRegex.FindAllStringSubmatch(line, -1)
	if len(matches) > 0 {
		for _, matche := range matches {
			title, destination := matche[1], matche[2]
			if strings.HasPrefix(destination, "http") {
				continue
			}
			fileName := filepath.Base(destination)
			fileName = strings.ReplaceAll(fileName, ".md", "")
			if fileName == title {
				if files, ok := c.filemap[fileName]; ok && len(files) >= 2 {
					file := strings.TrimPrefix(destination, "./")
					file = strings.TrimSuffix(file, ".md")
					line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s|%s]]`, file, title), 1)
				} else {
					line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s]]`, title), 1)
				}
			} else {
				file := strings.TrimPrefix(destination, "./")
				file = strings.TrimSuffix(file, ".md")
				line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s|%s]]`, file, title), 1)
			}
		}
	}
	return line
}
