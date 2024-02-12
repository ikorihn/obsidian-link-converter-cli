package olconv

import (
	"bytes"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	r := strings.NewReader(
		`---
title: My Note
tags:
  - Go
---

Hello convert

## How to use

https://github.com/ikorihn

- [this is link](https://example.com)
- [this is note](./Note.md)

[Note](./Note.md) is not [the link](http://example.com)

[paren in (link)](./Paren.md)

[hoge](./hoge.tar.md)
[samename](./sub1/samename.md)
[samename](./sub2/samename.md)
` +
			"```\n" +
			"[in codeblock](./codeblock.md)\n" +
			"```\n")

	w := &bytes.Buffer{}

	c := NewConverter(
		map[string][]string{
			"samename": {"./sub1/samename.md", "./sub2/samename.md"},
			"hoge":     {"./hoge.tar.md"},
		},
	)
	err := c.Convert(r, w)
	if err != nil {
		t.Errorf("Convert() error = %v", err)
		return
	}

	want := `---
title: My Note
tags:
  - Go
---

Hello convert

## How to use

https://github.com/ikorihn

- [this is link](https://example.com)
- [[Note|this is note]]

[[Note]] is not [the link](http://example.com)

[[Paren|paren in (link)]]

[[hoge.tar|hoge]]
[[sub1/samename|samename]]
[[sub2/samename|samename]]
` +
		"```\n" +
		"[in codeblock](./codeblock.md)\n" +
		"```\n"

	if gotW := w.String(); gotW != want {
		t.Errorf("Convert() = %v, want %v", gotW, want)
	}
}
