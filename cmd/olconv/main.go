package main

import (
	"flag"

	"github.com/ikorihn/olconv"
)

func main() {
	var basepath string
	var reverse bool
	flag.StringVar(&basepath, "basepath", ".", "specify target directory")
	flag.BoolVar(&reverse, "reverse", false, "convert Wikilinks to Markdown links (default: Markdown links to Wikilinks)")
	flag.Parse()

	var err error
	if reverse {
		err = olconv.ReverseConvertUnderDir(basepath)
	} else {
		err = olconv.ConvertUnderDir(basepath)
	}

	if err != nil {
		panic(err)
	}
}
