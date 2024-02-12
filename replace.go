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

func Convert(r io.Reader, w io.Writer) error {

	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)

	// file一覧をmapで保持
	// 同一ファイル名のときに判別するため、map[string][]string にする
	// codeblockは無視する

	inCodeBlock := false
	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
		}
		if inCodeBlock {
			bw.WriteString(line)
			bw.WriteString("\n")
			continue
		}

		matches := mdLinkRegex.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			for i, matche := range matches {
				title, destination := matche[1], matche[2]
				if strings.HasPrefix(destination, "http") {
					continue
				}
				fmt.Printf("Match! %d %v ==== %v\n", i, title, destination)
				fileName := filepath.Base(destination)
				fileName = strings.ReplaceAll(fileName, ".md", "")
				if fileName == title {
					line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s]]`, title), 1)
				} else {
					file := strings.TrimPrefix(destination, "./")
					file = strings.TrimSuffix(file, ".md")
					line = strings.Replace(line, fmt.Sprintf(`[%s](%s)`, title, destination), fmt.Sprintf(`[[%s|%s]]`, file, title), 1)
				}
			}
		}
		bw.WriteString(line)
		bw.WriteString("\n")
	}
	bw.Flush()

	return nil
}
