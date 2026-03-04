package cli

import (
	"errors"
	"flag"
	"strings"

	"yt-grab/internal/config"
	"yt-grab/internal/runner"
)

func runGrab(cfg config.Config, args []string) error {
	fs := flag.NewFlagSet("grab", flag.ContinueOnError)
	output := fs.String("output", "", "Output directory")
	format := fs.String("format", "", "yt-dlp format selector")
	maxRes := fs.Int("max-res", 0, "Maximum video resolution height")
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) != 1 {
		return errors.New("usage: yt-grab grab <url>")
	}
	url := rest[0]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("url must start with http:// or https://")
	}
	return runner.Run(runner.Request{URL: url, OutputDir: choose(*output, cfg.OutputDir), Format: choose(*format, cfg.Format), MaxResolution: chooseInt(*maxRes, cfg.MaxResolution), YtDLPPath: cfg.YtDLPPath, NoPlaylist: true})
}

func choose(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
func chooseInt(v, fallback int) int {
	if v != 0 {
		return v
	}
	return fallback
}
