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
[Note](./Note.md) and [Second Note](./Second%20Note.md)

[paren in (link)](./Paren.md)

[only bracket]
(only paren)
[not link] (only paren)

[hoge](./hoge.tar.md)
[samename](./sub1/samename.md)
[samename](./sub2/samename.md)
[日本語](./日本語.md)
` +
			"in `[code span](./codespan.md)`\n" +
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
	err := c.Convert(r, w, true)
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
[[Note]] and [[Second Note]]

[[Paren|paren in (link)]]

[only bracket]
(only paren)
[not link] (only paren)

[[hoge.tar|hoge]]
[[sub1/samename|samename]]
[[sub2/samename|samename]]
[[日本語]]
` +
		"in `[code span](./codespan.md)`\n" +
		"```\n" +
		"[in codeblock](./codeblock.md)\n" +
		"```\n"

	if gotW := w.String(); gotW != want {
		t.Errorf("Convert() = %v, want %v", gotW, want)
	}
}

func TestConverter_convertLine(t *testing.T) {
	type fields struct {
		inCodeBlock bool
		filemap     map[string][]string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{

		{
			name: "",
			fields: fields{
				filemap: map[string][]string{
					"circuit breaker pattern": {"note/circuit breaker pattern.md"},
				},
			},
			args: args{
				line: `[circuit breaker](note/circuit%20breaker%20pattern.md) を実装する`,
			},
			want: "[[circuit breaker pattern|circuit breaker]] を実装する",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Converter{
				inCodeBlock: tt.fields.inCodeBlock,
				filemap:     tt.fields.filemap,
			}
			if got := c.convertLine(tt.args.line); got != tt.want {
				t.Errorf("Converter.convertLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_readChar(t *testing.T) {
	l := Parser{
		mdLinks: []mdLink{},
	}
	l.parse("[circuit breaker](note/circuit%20breaker%20pattern.md) を実装する`")

	t.Logf("%#v\n", l)
}
