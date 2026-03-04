# yt-grab

yt-grab is a fast, cross-platform media downloader and automation-oriented CLI written in Go.
It delegates extraction to `yt-dlp` and focuses on a clean command interface, reliable orchestration, and scriptability.

## Features

- `grab`, `audio`, and `playlist` subcommands
- Config support via `~/.ytgrab.yaml`
- Dependency diagnostics via `doctor`
- Cross-platform output path defaults
- Explicit wrapper around `yt-dlp` (no scraping reimplementation)

## Installation

### Prerequisites
- Go 1.22+
- `yt-dlp` in `PATH`
- `ffmpeg` in `PATH` (recommended for merges/conversion)

### Build

```bash
go build -o yt-grab ./cmd/yt-grab
```

## Usage

```bash
yt-grab [--config ~/.ytgrab.yaml] <command> [flags]
```

### Commands

- `yt-grab grab <url> [--output DIR] [--format FORMAT] [--max-res 1080]`
- `yt-grab audio <url> [--output DIR] [--format mp3|aac|wav]`
- `yt-grab playlist <url> [--output DIR] [--flat]`
- `yt-grab config init|view`
- `yt-grab doctor`
- `yt-grab version`

## Config

Create defaults:

```bash
yt-grab config init
```

Default config path: `~/.ytgrab.yaml`

Example:

```yaml
output_dir: ~/Downloads
format: bestvideo+bestaudio/best
audio_format: mp3
max_resolution: 1080
yt_dlp_path: ""
```

Config precedence:

`command flags > env vars > config file > defaults`

## Notes

- yt-grab does not bypass DRM/paywalls.
- Users are responsible for complying with local laws and platform terms.
