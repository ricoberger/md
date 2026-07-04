# md

`md` is a simple Markdown renderer for the terminal that uses
[Glamour](https://github.com/charmbracelet/glamour).

It exists because Glow didn't work well within the
[fzf](https://github.com/junegunn/fzf) preview when using custom styles and
colors. `md` honors the `GLAMOUR_STYLE` environment variable and works with the
[Catppuccin Glamour theme](https://github.com/catppuccin/glamour).

## Install

### Releases

Check the [releases](https://github.com/ricoberger/md/releases) page for the
full list of pre-built binaries.

1. Download the release for your os/arch
2. Unzip the archive to get the `md` binary
3. Add the `md` binary to your `PATH`

### Source

```sh
go install github.com/ricoberger/md@latest
```

## Development

To build and run the binary the following commands can be used:

```sh
go build -o ./bin/md .
./bin/md
```
