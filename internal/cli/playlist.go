package cli

import (
	"errors"
	"flag"
	"strings"

	"yt-grab/internal/config"
	"yt-grab/internal/runner"
)

func runPlaylist(cfg config.Config, args []string, plainOutput bool) error {
	fs := flag.NewFlagSet("playlist", flag.ContinueOnError)
	output := fs.String("output", "", "Output directory")
	flat := fs.Bool("flat", false, "Use --flat-playlist metadata mode")
	plain := fs.Bool("plain-output", plainOutput, "Use raw yt-dlp output instead of compact progress bar")
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) != 1 {
		return errors.New("usage: yt-grab playlist <url>")
	}
	url := rest[0]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("url must start with http:// or https://")
	}
	return runner.Run(runner.Request{URL: url, OutputDir: choose(*output, cfg.OutputDir), Format: cfg.Format, PlainOutput: *plain, FlatPlaylist: *flat, YtDLPPath: cfg.YtDLPPath})
}
