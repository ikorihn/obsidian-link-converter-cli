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
	inCodeSpan          bool
	inMdLinkTitle       bool
	inMdLinkDestination bool

	mdLinks []mdLink
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
	var mdLinkTitle, mdLinkDestination string
	var titleStartPos, titleEndPos, destinationStartPos, destinationEndPos int

	var skipNext bool

	for i, c := range input {
		// XXX ugly code ...
		if skipNext {
			skipNext = false
			continue
		}

		if c == '[' && !p.inCodeSpan && !p.inMdLinkTitle && !p.inMdLinkDestination {
			p.inMdLinkTitle = true
			titleStartPos = i
			mdLinkTitle = ""

		} else if c == ']' && !p.inCodeSpan && p.inMdLinkTitle && !p.inMdLinkDestination {
			if i+1 < len(input) && input[i+1] == '(' {
				// only if next char is "("
				titleEndPos = i
				skipNext = true
				p.inMdLinkDestination = true
				destinationStartPos = i
				mdLinkDestination = ""
			}

			p.inMdLinkTitle = false

		} else if c == ')' && !p.inCodeSpan && !p.inMdLinkTitle && p.inMdLinkDestination {

			p.inMdLinkDestination = false
			destinationEndPos = i

			p.mdLinks = append(p.mdLinks, mdLink{
				title:               mdLinkTitle,
				destination:         mdLinkDestination,
				titleStartPos:       titleStartPos,
				titleEndPos:         titleEndPos,
				destinationStartPos: destinationStartPos,
				destinationEndPos:   destinationEndPos,
			})

		} else if c == '`' {

			if p.inMdLinkTitle {
				mdLinkTitle += string(c)
			} else if p.inMdLinkDestination {
				mdLinkDestination += string(c)
			} else {
				p.inCodeSpan = !p.inCodeSpan
			}

		} else {

			if p.inMdLinkTitle {
				mdLinkTitle += string(c)
			} else if p.inMdLinkDestination {
				mdLinkDestination += string(c)
			}

		}
	}

}
