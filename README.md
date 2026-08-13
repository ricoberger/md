# md

`md` is a simple Markdown renderer for the terminal that uses
[Glamour](https://github.com/charmbracelet/glamour).

It exists because Glow didn't work well within the
[fzf](https://github.com/junegunn/fzf) preview when using custom styles and
colors. `md` honors the `GLAMOUR_STYLE` environment variable and works with the
[Catppuccin Glamour theme](https://github.com/catppuccin/glamour).

- `MD_WORD_WRAP`: Set the word wrap width. Defaults to `80`.
- `MD_ADF`: Set to `true` if the input is in Atlassian Document Format (ADF).
  Defaults to `false`.
- `MD_RAW`: Set to `true` to print the raw input instead of rendering it.
  Defaults to `false`.

A leading YAML front matter block (delimited by `---`) is stripped before
rendering, so metadata such as tags is not printed. When `MD_RAW` is set the
input is printed verbatim, including any front matter.

- `MD_ENABLE_FRONT_MATTER`: Set to `true` to keep (render) the YAML front
  matter instead of stripping it. Defaults to `false`.

## Install

### Releases

Check the [releases](https://github.com/ricoberger/md/releases) page for the
full list of pre-built binaries.

1. Download the release for your os/arch
2. Unzip the archive to get the `md` binary
3. Add the `md` binary to your `PATH`

### Source

```sh
GOEXPERIMENT=jsonv2 go install github.com/ricoberger/md@latest
```

## Development

To build and run the binary the following commands can be used:

```sh
GOEXPERIMENT=jsonv2 go build -o ./bin/md .
./bin/md
```
