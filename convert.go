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

	lines := make([]string, 0)
	for sc.Scan() {
		line := sc.Text()
		line = c.convertLine(line)
		lines = append(lines, line)
	}
	bw.WriteString(strings.Join(lines, "\n"))
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

// ReverseConverter converts Wikilinks to Markdown links
type ReverseConverter struct {
	inCodeBlock bool
	filemap     map[string][]string
}

func NewReverseConverter(filemap map[string][]string) *ReverseConverter {
	return &ReverseConverter{
		filemap: filemap,
	}
}

func (rc *ReverseConverter) Convert(r io.Reader, w io.Writer, newLineAtEnd bool) error {
	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)

	lines := make([]string, 0)
	for sc.Scan() {
		line := sc.Text()
		line = rc.convertLine(line)
		lines = append(lines, line)
	}
	bw.WriteString(strings.Join(lines, "\n"))
	if newLineAtEnd {
		bw.WriteString("\n")
	}

	rc.inCodeBlock = false
	bw.Flush()

	return nil
}

func (rc *ReverseConverter) convertLine(line string) string {
	if strings.HasPrefix(line, "```") {
		rc.inCodeBlock = !rc.inCodeBlock
	}
	if rc.inCodeBlock {
		return line
	}

	wp := WikilinkParser{
		wikilinks: []wikilink{},
	}

	wp.parse(line)

	// start from last index to avoid index misalignment due to re-slicing
	for i := len(wp.wikilinks) - 1; i >= 0; i-- {
		wlink := wp.wikilinks[i]

		var mdLink string
		if wlink.title != "" {
			// [[destination|title]] -> [title](destination.md)
			mdLink = fmt.Sprintf(`[%s](%s.md)`, wlink.title, wlink.destination)
		} else {
			// [[destination]] -> [destination](destination.md)
			filename := extractFilename(wlink.destination)
			mdLink = fmt.Sprintf(`[%s](%s.md)`, filename, wlink.destination)
		}

		line = line[:wlink.startPos] + mdLink + line[wlink.endPos+1:]
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

// WikilinkParser parses WikiLinks in text
type WikilinkParser struct {
	inCodeSpan bool
	wikilinks  []wikilink
}

type wikilink struct {
	destination string
	title       string
	startPos    int
	endPos      int
}

func (wp *WikilinkParser) parse(input string) {
	var currentWikilink *wikilink
	i := 0

	for i < len(input) {
		switch input[i] {
		case '`':
			wp.inCodeSpan = !wp.inCodeSpan
		case '[':
			if wp.inCodeSpan {
				i++
				continue
			}
			// Check for [[
			if i+1 < len(input) && input[i+1] == '[' {
				currentWikilink = &wikilink{
					startPos: i,
				}
				i += 2 // Skip [[
				continue
			}
		case ']':
			if wp.inCodeSpan || currentWikilink == nil {
				i++
				continue
			}
			// Check for ]]
			if i+1 < len(input) && input[i+1] == ']' {
				// Extract content between [[ and ]]
				content := input[currentWikilink.startPos+2 : i]
				currentWikilink.endPos = i + 1

				// Parse content: check for | separator
				if pipeIndex := strings.Index(content, "|"); pipeIndex != -1 {
					currentWikilink.destination = strings.TrimSpace(content[:pipeIndex])
					currentWikilink.title = strings.TrimSpace(content[pipeIndex+1:])
				} else {
					currentWikilink.destination = strings.TrimSpace(content)
				}

				wp.wikilinks = append(wp.wikilinks, *currentWikilink)
				currentWikilink = nil
				i += 2 // Skip ]]
				continue
			}
		}
		i++
	}
}

// extractFilename extracts the filename from a path
func extractFilename(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
