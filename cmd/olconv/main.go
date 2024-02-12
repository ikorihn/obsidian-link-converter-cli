package main

import (
	"flag"

	"github.com/ikorihn/olconv"
)

func main() {
	var basepath string
	flag.StringVar(&basepath, "basepath", ".", "specify target directory")
	flag.Parse()

	if err := olconv.ConvertUnderDir(basepath); err != nil {
		panic(err)
	}

}
