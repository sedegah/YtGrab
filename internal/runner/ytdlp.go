package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Request struct {
	URL           string
	OutputDir     string
	Format        string
	AudioOnly     bool
	AudioFormat   string
	MaxResolution int
	NoPlaylist    bool
	FlatPlaylist  bool
	YtDLPPath     string
}

func BuildArgs(req Request) []string {
	args := []string{
		"--no-part",
		"--progress",
		"--progress-template",
		"download:%(progress._percent_str)s at %(progress._speed_str)s ETA %(progress._eta_str)s",
		"-o",
		filepath.Join(req.OutputDir, "%(title)s.%(ext)s"),
	}

	if req.NoPlaylist {
		args = append(args, "--no-playlist")
	}
	if req.FlatPlaylist {
		args = append(args, "--flat-playlist")
	}

	if req.AudioOnly {
		args = append(args, "-x", "--audio-format", req.AudioFormat)
	} else {
		format := req.Format
		if req.MaxResolution > 0 {
			cap := strconv.Itoa(req.MaxResolution)
			format = "bestvideo[height<=" + cap + "]+bestaudio/best[height<=" + cap + "]"
		}
		if format == "" {
			format = "bestvideo+bestaudio/best"
		}
		args = append(args, "-f", format)
	}

	return append(args, req.URL)
}

func Run(req Request) error {
	bin := req.YtDLPPath
	if bin == "" {
		resolved, err := exec.LookPath("yt-dlp")
		if err != nil {
			return fmt.Errorf("yt-dlp not found in PATH: %w", err)
		}
		bin = resolved
	}

	cmd := exec.Command(bin, BuildArgs(req)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CheckDependency(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s not found in PATH", name)
	}
	return nil
}
