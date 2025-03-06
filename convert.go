package olconv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Converter struct {
	inCodeBlock bool
	filemap     map[string][]string
}

func NewConverter(filemap map[string][]string) *Converter {
	return &Converter{
		filemap: filemap,
	}
}

func (c *Converter) Convert(r io.Reader, w io.Writer, newLineAtEnd bool) error {
	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)

	for sc.Scan() {
		line := sc.Text()
		line = c.convertLine(line)
		bw.WriteString(line)
		bw.WriteString("\n")
	}
	c.inCodeBlock = false
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

	p := Parser{
		mdLinks: []mdLink{},
	}

	p.parse(line)

	// start from last index to avoid index misalignment due to re-slicing
	for i := len(p.mdLinks) - 1; i >= 0; i-- {
		mdLink := p.mdLinks[i]
		title, destination := mdLink.title, mdLink.destination
		if strings.HasPrefix(destination, "http") {
			continue
		}

		filename := filenameWithoutMdExtension(destination)

		if filename == title {
			if files, ok := c.filemap[filename]; ok && len(files) >= 2 {
				relativePath := formatRelativePath(destination)
				line = line[:mdLink.titleStartPos] + fmt.Sprintf(`[[%s|%s]]`, relativePath, title) + line[mdLink.destinationEndPos+1:]
			} else {
				line = line[:mdLink.titleStartPos] + fmt.Sprintf(`[[%s]]`, title) + line[mdLink.destinationEndPos+1:]
			}
		} else {
			if files, ok := c.filemap[filename]; ok && len(files) == 1 {
				line = line[:mdLink.titleStartPos] + fmt.Sprintf(`[[%s|%s]]`, filename, title) + line[mdLink.destinationEndPos+1:]
			} else {
				relativePath := formatRelativePath(destination)
				line = line[:mdLink.titleStartPos] + fmt.Sprintf(`[[%s|%s]]`, relativePath, title) + line[mdLink.destinationEndPos+1:]
			}
		}

	}

	return line
}

type Parser struct {
	inCodeSpan bool
	mdLinks    []mdLink
}

type mdLink struct {
	title               string
	destination         string
	titleStartPos       int
	titleEndPos         int
	destinationStartPos int
	destinationEndPos   int
}

func (p *Parser) parse(input string) {
	var currentLink *mdLink
	for i, c := range input {
		switch c {
		case '`':
			p.inCodeSpan = !p.inCodeSpan
		case '[':
			if p.inCodeSpan {
				continue
			}
			currentLink = &mdLink{
				titleStartPos: i,
			}
		case ']':
			if p.inCodeSpan || currentLink == nil {
				continue
			}
			currentLink.titleEndPos = i
			if currentLink.destinationStartPos == 0 {
				currentLink.title = input[currentLink.titleStartPos+1 : i]
			}
			if i+1 < len(input) && input[i+1] == '(' {
				currentLink.destinationStartPos = i + 1
			} else {
				currentLink = nil
			}
		case '(':
			if p.inCodeSpan || currentLink == nil {
				continue
			}
			if currentLink.titleEndPos == 0 {
				continue
			}
		case ')':
			if p.inCodeSpan || currentLink == nil || currentLink.destinationStartPos == 0 {
				continue
			}

			currentLink.destination = input[currentLink.destinationStartPos+1 : i]
			currentLink.destinationEndPos = i
			p.mdLinks = append(p.mdLinks, *currentLink)
			currentLink = nil
		default:
			continue
		}
	}
}
