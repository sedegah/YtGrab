package downloader

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Options struct {
	URL        string
	OutputDir  string
	AudioOnly  bool
	AudioFmt   string
	MaxHeight  int
	Format     string
	NoPlaylist bool
}

type Runner struct {
	path   string
	stdout io.Writer
	stderr io.Writer
}

func NewRunner(path string, stdout, stderr io.Writer) (*Runner, error) {
	if path == "" {
		resolved, err := exec.LookPath("yt-dlp")
		if err != nil {
			return nil, fmt.Errorf("yt-dlp not found in PATH; install it or pass --yt-dlp: %w", err)
		}
		path = resolved
	}
	return &Runner{path: path, stdout: stdout, stderr: stderr}, nil
}

func BuildArgs(opts Options) []string {
	args := []string{
		"--newline",
		"--no-part",
		"-o", filepath.Join(opts.OutputDir, "%(title)s.%(ext)s"),
	}

	if opts.NoPlaylist {
		args = append(args, "--no-playlist")
	}

	if opts.AudioOnly {
		args = append(args,
			"-x",
			"--audio-format", opts.AudioFmt,
		)
	} else {
		if opts.Format != "" {
			args = append(args, "-f", opts.Format)
		} else if opts.MaxHeight > 0 {
			args = append(args, "-f", "bestvideo[height<="+strconv.Itoa(opts.MaxHeight)+"]+bestaudio/best[height<="+strconv.Itoa(opts.MaxHeight)+"]")
		} else {
			args = append(args, "-f", "bestvideo+bestaudio/best")
		}
	}

	args = append(args, opts.URL)
	return args
}

func (r *Runner) Run(args []string) error {
	cmd := exec.Command(r.path, args...)
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
