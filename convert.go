package olconv

import (
	"bufio"
	"fmt"
	"io"
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

func (c *Converter) Convert(r io.Reader, w io.Writer, newLineAtEnd bool) error {

	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)

	first := true
	for sc.Scan() {
		line := sc.Text()

		line = c.convertLine(line)
		if !first {
			bw.WriteString("\n")
		}
		bw.WriteString(line)
		if first {
			first = false
		}
	}
	if newLineAtEnd {
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

			filename := filenameWithoutMdExtension(destination)

			if filename == title {
				if files, ok := c.filemap[filename]; ok && len(files) >= 2 {
					relativePath := formatRelativePath(destination)
					line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s|%s]]`, relativePath, title), 1)
				} else {
					line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s]]`, title), 1)
				}
			} else {
				relativePath := formatRelativePath(destination)
				line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s|%s]]`, relativePath, title), 1)
			}
		}
	}
	return line
}
