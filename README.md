# Obsidian Link Converter CLI

**Obsidian Link Converter CLI** is a command-line tool that converts `Markdown Link` within an Obsidian vault to `WikiLinks` link format.

### Features

This cli converts `Markdown Link` within an Obsidian vault to `WikiLinks` link format.
The cli converts final path to shortest possible.

### Usage

```shell
❯ olconv [-basepath path/to/your/vault]
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
