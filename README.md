# Obsidian Link Converter CLI

![alt text](https://img.shields.io/badge/cli-blue)
![alt text](https://img.shields.io/badge/Go-00ADD8?logo=go)

**Obsidian Link Converter CLI** is a command-line tool that converts `Markdown Link` within an Obsidian vault to `Wikilink` link format and vice versa.

## Features

This CLI tool converts Markdown links within your Obsidian vault to Wikilink format.

### Conversion Examples

| Description              | Markdown link                                  | Wikilink                                 |
| ------------------------ | ---------------------------------------------- | ---------------------------------------- |
| Basic conversion         | `[note](note.md)`                              | `[[note]]`                               |
| Link with title          | `[Title](note.md)`                             | `[[note\|Title]]`                        |
| File in subfolder        | `[note](subfolder/note.md)`                    | `[[subfolder/note\|note]]`               |
| Same filename detection  | `[foo](./sub1/foo.md)`, `[foo](./sub2/foo.md)` | `[[sub1/foo\|foo]]`, `[[sub2/foo\|foo]]` |
| Shortest path conversion | `[foo](./path/to/subdir/foo.md)`               | `[[foo]]`                                |

### Additional Features

- Does not convert links in code blocks
- Does not convert links in code spans
- Does not convert external links
- Converts to shortest path possible

## Note

- When converting links, please be aware that the shortest path conversion may result in path information loss during multiple conversions. For example:
`[foo](./path/to/subdir/foo.md)` -- `-to-wiki` --> `[[foo]]` -- `-to-markdown` --> `[foo](foo.md)`

- It's recommended to create a backup of your Obsidian vault before using this tool.
- Verify that the converted Markdown files work correctly.


## Installation

```shell
❯ go install github.com/ikorihn/olconv/cmd/olconv@latest
```

## Usage

```shell
# Convert Markdown links to Wikilink in the current directory
❯ olconv -to-wiki

# Convert Wikilink to Markdown links in the current directory
❯ olconv -to-markdown

# Convert links in a specific directory
❯ olconv -basepath path/to/your/vault -to-wiki
❯ olconv -basepath path/to/your/vault -to-markdown

# Show help
❯ olconv -h
```

### Command Line Options

- `-to-wiki`: Convert Markdown links `[title](path.md)` to Wikilink `[[path|title]]`
- `-to-markdown`: Convert Wikilink `[[path|title]]` to Markdown links `[title](path.md)`
- `-basepath <path>`: Specify target directory (default: current directory)

**Note**: You must specify either `-to-wiki` or `-to-markdown` (but not both).

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

### License

This project is licensed under the [MIT License](LICENSE).
