package main

import (
	"os"

	"github.com/ikorihn/olconv"
)

func main() {
	f, err := os.OpenFile("test.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	wf, err := os.Create("after.md")
	if err != nil {
		panic(err)
	}
	defer wf.Close()

	if err := olconv.Convert(f, wf); err != nil {
		panic(err)
	}
}
