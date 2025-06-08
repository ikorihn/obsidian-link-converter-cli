package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikorihn/olconv"
)

func main() {
	var basepath string
	var toWiki bool
	var toMarkdown bool

	flag.StringVar(&basepath, "basepath", ".", "specify target directory")
	flag.BoolVar(&toWiki, "to-wiki", false, "convert Markdown links to Wikilinks")
	flag.BoolVar(&toMarkdown, "to-markdown", false, "convert Wikilinks to Markdown links")
	flag.Parse()

	// どちらも指定されていない、または両方指定されている場合
	if (!toWiki && !toMarkdown) || (toWiki && toMarkdown) {
		fmt.Fprintf(os.Stderr, "Error: Please specify either --to-wiki or --to-markdown\n\n")
		flag.Usage()
		os.Exit(1)
	}

	var err error
	if toWiki {
		err = olconv.LinkToWikilink(basepath)
	} else {
		err = olconv.WikilinkToLink(basepath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
