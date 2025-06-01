# Obsidian Link Converter CLI

![alt text](https://img.shields.io/badge/cli-blue)
![alt text](https://img.shields.io/badge/Go-00ADD8?logo=go)

**Obsidian Link Converter CLI** is a command-line tool that converts `Markdown Link` within an Obsidian vault to `WikiLinks` link format.

### Features

This CLI tool converts Markdown links within your Obsidian vault to WikiLinks format.

- Converts `[Example](note.md)` to `[[note]]`
- Converts `[Title](./note.md)` to `[[note|Title]]`
- Converts `[Example](subfolder/note.md)` to `[[subfolder/note]]`
- automatically detect same file name in different folder, so convert `[foo](./sub1/foo.md)` and `[foo](./sub2/foo.md)` to `[[sub1/foo|foo]]` and `[[sub2/foo|foo]]`
- Removes `.md` from the link target.
- Handles links containing spaces and special characters.
- Handles codeblock. Links in codeblock will not convert.
- Handles code span. Links in code span will not convert.
- Handles external links. External links will not convert.
- Removes `/` at the first of the link path.
- Removes `.md` at the end of the link path.
- Converts to shortest path possible. for example if you have `./sub1/foo.md`, `./sub2/foo.md` and `./foo.md`, when you add `[foo](./foo.md)` in your note, it will convert to `[[foo]]`.

### Usage

```shell
# Convert links in the current directory
❯ olconv
# Convert links in a specific directory
❯ olconv -basepath path/to/your/vault
# Show help
❯ olconv -h
```

Input

```
- [This is Link](note.md)
- [This is Link](subfolder/note.md)
```

Output

```
- [[note|This is Link]]
- [[subfolder/note|This is Link]]
```

### Background

The reason for creating this tool was the significant time it took to process link conversion using plugins. Therefore, by writing this CLI in Go, link conversion can now be completed almost instantly.

### Notes

- It's recommended to create a backup of your Obsidian vault before using this tool.
- Verify that the converted Markdown files work correctly.

### Installation

```shell
❯ go install github.com/ikorihn/obsidian-link-converter-cli/cmd/olconv@latest
```

### License

This project is licensed under the [MIT License](LICENSE).
