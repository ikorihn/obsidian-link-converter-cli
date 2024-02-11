package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	f, err := os.Open("test.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	wf, err := os.Create("after.md")
	if err != nil {
		panic(err)
	}
	defer wf.Close()

	if err := run(f, wf); err != nil {
		panic(err)
	}
}

var mdLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)

func run(r io.Reader, w io.Writer) error {

	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)

	// file一覧をmapで保持
	// 同一ファイル名のときに判別するため、map[string][]string にする
	// codeblockは無視する

	for sc.Scan() {
		line := sc.Text()
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

func runMd() error {
	source, err := os.ReadFile("test.md")
	if err != nil {
		return err
	}

	gm := goldmark.New()
	n := gm.Parser().Parse(text.NewReader(source))
	n.Dump(source, 2)

	err = ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {

		switch kind := n.(type) {
		case *ast.FencedCodeBlock, *ast.CodeSpan, *ast.CodeBlock:
			return ast.WalkContinue, nil
		case *ast.Link:
			if strings.HasPrefix(string(kind.Destination), "http") {
				return ast.WalkSkipChildren, nil
			}
			fmt.Printf("Link: [%s](%s)\n", string(kind.Title), string(kind.Destination))
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})

	return err
}
