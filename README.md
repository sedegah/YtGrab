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

## Django Frontend

A Django landing page is available in [frontend/](frontend/) and styled to match the referenced hero-card design.
It also includes a working download form:

- Paste just the link for video download
- Paste `link --audio` to download audio

### Run locally

```bash
cd frontend
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python manage.py migrate
python manage.py runserver
```

Open `http://127.0.0.1:8000/`.

### Windows (PowerShell)

PowerShell does not run executables from the current folder unless you prefix with `./` (or `\.\`).

```powershell
go build -o yt-grab.exe ./cmd/yt-grab
.\yt-grab.exe --help
```

To use `yt-grab` directly without `\.\`, install it to a folder in `PATH` (for example `$env:USERPROFILE\go\bin`).

Quick installer (build + PATH setup):

```powershell
.\install.ps1
```

Install with dependencies (`yt-dlp` and `ffmpeg`) using winget:

```powershell
.\install.ps1 -InstallDeps
```

Optional custom install directory:

```powershell
.\install.ps1 -InstallDir "$env:USERPROFILE\bin"
```

Windows default download folder is:

```text
C:\Users\<YourUser>\Downloads\YtGrab
```

Example:

```powershell
.\yt-grab.exe https://youtu.be/8ekJMC8OtGU --output "$env:USERPROFILE\Downloads\YtGrab"
```

If dependencies are missing:

```powershell
winget install yt-dlp.yt-dlp
winget install Gyan.FFmpeg
```

## Usage

```bash
yt-grab [--config ~/.ytgrab.yaml] [--plain-output] <url> [--audio] [--quality best|worst|720p|1080p] [--output DIR]
```

Also supported:

```bash
yt-grab [--config ~/.ytgrab.yaml] <command> [flags]
```

### Quick examples

- `yt-grab https://youtu.be/8ekJMC8OtGU`
- `yt-grab https://youtu.be/8ekJMC8OtGU --audio`
- `yt-grab https://youtu.be/8ekJMC8OtGU --quality 1080p`
- `yt-grab https://youtu.be/8ekJMC8OtGU --quality best --output ~/Downloads`
- `yt-grab --plain-output https://youtu.be/8ekJMC8OtGU`

### Commands

- `yt-grab grab <url> [--output DIR] [--format FORMAT] [--max-res 1080] [--quality best|worst|720p|1080p]`
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

Windows example:

```yaml
output_dir: C:\Users\YourUser\Downloads\YtGrab
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
