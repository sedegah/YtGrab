package cli

import (
	"errors"
	"flag"
	"strings"

	"yt-grab/internal/config"
	"yt-grab/internal/runner"
)

func runGrab(cfg config.Config, args []string, plainOutput bool) error {
	fs := flag.NewFlagSet("grab", flag.ContinueOnError)
	output := fs.String("output", "", "Output directory")
	format := fs.String("format", "", "yt-dlp format selector")
	maxRes := fs.Int("max-res", 0, "Maximum video resolution height")
	quality := fs.String("quality", "", "Video quality: best|worst|720p|1080p")
	plain := fs.Bool("plain-output", plainOutput, "Use raw yt-dlp output instead of compact progress bar")
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
	selectedFormat := choose(*format, cfg.Format)
	selectedMaxRes := chooseInt(*maxRes, cfg.MaxResolution)
	if *quality != "" {
		if *format != "" {
			return errors.New("use either --format or --quality, not both")
		}
		qFormat, qMaxRes, err := qualityToFormat(*quality)
		if err != nil {
			return err
		}
		if qFormat != "" {
			selectedFormat = qFormat
			selectedMaxRes = 0
		}
		if qMaxRes > 0 {
			selectedMaxRes = qMaxRes
		}
	}
	return runner.Run(runner.Request{URL: url, OutputDir: choose(*output, cfg.OutputDir), Format: selectedFormat, PlainOutput: *plain, MaxResolution: selectedMaxRes, YtDLPPath: cfg.YtDLPPath, NoPlaylist: true})
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
