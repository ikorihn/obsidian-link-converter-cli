package olconv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			"```go\n" +
			"[in codeblock](./codeblock.md)\n" +
			"```\n" +
			"in `[code span](./codespan.md)`\n")

	w := &bytes.Buffer{}

	c := NewConverter(
		map[string][]string{
			"samename": {"./sub1/samename.md", "./sub2/samename.md"},
			"hoge":     {"./hoge.tar.md"},
		},
	)
	err := c.Convert(r, w, true, ToWikilink)
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
		"```go\n" +
		"[in codeblock](./codeblock.md)\n" +
		"```\n" +
		"in `[code span](./codespan.md)`\n"

	assert.Equal(t, strings.Split(want, "\n"), strings.Split(w.String(), "\n"))
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
					"note with spaces": {"note/note with spaces.md"},
				},
			},
			args: args{
				line: `[note title](note/note%20with%20spaces.md) を実装する`,
			},
			want: "[[note with spaces|note title]] を実装する",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Converter{
				inCodeBlock: tt.fields.inCodeBlock,
				filemap:     tt.fields.filemap,
			}
			if got := c.convertLine(tt.args.line, ToWikilink); got != tt.want {
				t.Errorf("Converter.convertLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_readChar(t *testing.T) {
	l := Parser{
		mdLinks: []mdLink{},
	}
	l.parse("[note title](note/note%20with%20spaces.md) を実装する`")

	t.Logf("%#v\n", l)
}

func TestReverseConvert(t *testing.T) {
	r := strings.NewReader(
		`---
title: My Note
tags:
  - Go
---

Hello reverse convert

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
			"in `[[codespan]]`\n" +
			"```go\n" +
			"[[in codeblock]]\n" +
			"```\n" +
			"in `[[codespan]]`\n")

	w := &bytes.Buffer{}

	c := NewConverter(
		map[string][]string{
			"samename": {"./sub1/samename.md", "./sub2/samename.md"},
			"hoge":     {"./hoge.tar.md"},
		},
	)
	err := c.Convert(r, w, true, ToMarkdown)
	if err != nil {
		t.Errorf("ReverseConvert() error = %v", err)
		return
	}

	want := `---
title: My Note
tags:
  - Go
---

Hello reverse convert

## How to use

https://github.com/ikorihn

- [this is link](https://example.com)
- [this is note](Note.md)

[Note](Note.md) is not [the link](http://example.com)
[Note](Note.md) and [Second Note](Second Note.md)

[paren in (link)](Paren.md)

[only bracket]
(only paren)
[not link] (only paren)

[hoge](hoge.tar.md)
[samename](sub1/samename.md)
[samename](sub2/samename.md)
[日本語](日本語.md)
` +
		"in `[[codespan]]`\n" +
		"```go\n" +
		"[[in codeblock]]\n" +
		"```\n" +
		"in `[[codespan]]`\n"

	assert.Equal(t, strings.Split(want, "\n"), strings.Split(w.String(), "\n"))
}

func TestReverseConverter_convertLine(t *testing.T) {
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
			name: "simple wikilink",
			fields: fields{
				filemap: map[string][]string{
					"Note": {"Note.md"},
				},
			},
			args: args{
				line: `This is a [[Note]] reference.`,
			},
			want: "This is a [Note](Note.md) reference.",
		},
		{
			name: "wikilink with title",
			fields: fields{
				filemap: map[string][]string{
					"note with spaces": {"note/note with spaces.md"},
				},
			},
			args: args{
				line: `[[note with spaces|note title]] を実装する`,
			},
			want: "[note title](note with spaces.md) を実装する",
		},
		{
			name: "wikilink with path",
			fields: fields{
				filemap: map[string][]string{
					"samename": {"./sub1/samename.md", "./sub2/samename.md"},
				},
			},
			args: args{
				line: `[[sub1/samename|samename]] and [[sub2/samename]]`,
			},
			want: "[samename](sub1/samename.md) and [samename](sub2/samename.md)",
		},
		{
			name: "multiple wikilinks",
			fields: fields{
				filemap: map[string][]string{
					"Note": {"Note.md"},
					"Other": {"Other.md"},
				},
			},
			args: args{
				line: `See [[Note]] and [[Other|another note]] for details.`,
			},
			want: "See [Note](Note.md) and [another note](Other.md) for details.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Converter{
				inCodeBlock: tt.fields.inCodeBlock,
				filemap:     tt.fields.filemap,
			}
			if got := c.convertLine(tt.args.line, ToMarkdown); got != tt.want {
				t.Errorf("Converter.convertLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWikilinkParser_parse(t *testing.T) {
	wp := WikilinkParser{
		wikilinks: []wikilink{},
	}
	wp.parse("[[note with spaces|note title]] を実装する")

	assert.Len(t, wp.wikilinks, 1)
	assert.Equal(t, "note with spaces", wp.wikilinks[0].destination)
	assert.Equal(t, "note title", wp.wikilinks[0].title)
	assert.Equal(t, 0, wp.wikilinks[0].startPos)
	assert.Equal(t, 30, wp.wikilinks[0].endPos)
}

func TestExtractFilename(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "simple filename",
			path: "Note",
			want: "Note",
		},
		{
			name: "path with directory",
			path: "sub1/samename",
			want: "samename",
		},
		{
			name: "nested path",
			path: "folder/subfolder/file",
			want: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractFilename(tt.path); got != tt.want {
				t.Errorf("extractFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
