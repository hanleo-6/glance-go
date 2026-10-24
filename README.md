# glance

`glance` is a Yet Another CLI Mark Down Parser also Hosting Directories,
or Yacli MD, PhD for short.

This version is written in Go.

## Features

- **GitHub Flavored Markdown (GFM):** Renders `.md` files the same way github does,
  via [goldmark](https://github.com/yuin/goldmark).
- **Live 'Live' Reload:** Watches `.md` files for changes like its life depends on it.
- **Relative Assets & Links:** Serves local images and relative links without
  snooping.
- **Directory Mode:** `glance` at a directory to browse `.md` files without switching tabs.
- **Intuitive Name:** It has an easy name to remember and type.

## Installation

Requires Go 1.22 or later:

```bash
go install github.com/hanleo-6/glance-go/cmd/glance@latest
```

Or build from source:

```bash
git clone https://github.com/hanleo-6/glance-go.git
cd glance-go
go build -o glance ./cmd/glance
```

## Usage

Preview a single Markdown file:

```bash
glance README.md
```

Preview a directory of Markdown documents:

```bash
glance ./docs
```

Specify a custom port:

```bash
glance -port 8080 README.md
```

## Mostly how it works

```text
editor save -> watcher -> debounce -> hub (broadcaster) -> SSE clients (browser tabs)
                                                                            |
browser <- HTML page (rendered Markdown + minimal JS) <--------------- HTTP server
```

### Project Layout

```text
glance-go/
├── cmd/
│   └── glance/       # CLI entrypoint and wiring
└── internal/
    ├── translate/    # Markdown to HTML conversion (using goldmark)
    ├── watch/        # Filesystem watcher
    ├── hub/          # Broadcast messages to its subscribers (browser tabs)
    └── server/       # HTTP server
```

## Development

Run directly during development:

```bash
go run ./cmd/glance README.md
```

Run tests (all zero of them):

```bash
go test ./...
```

