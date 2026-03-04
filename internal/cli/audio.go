package cli

import (
	"errors"
	"flag"
	"strings"

	"yt-grab/internal/config"
	"yt-grab/internal/runner"
)

func runAudio(cfg config.Config, args []string, plainOutput bool) error {
	fs := flag.NewFlagSet("audio", flag.ContinueOnError)
	output := fs.String("output", "", "Output directory")
	format := fs.String("format", "", "Audio format: mp3|aac|wav")
	plain := fs.Bool("plain-output", plainOutput, "Use raw yt-dlp output instead of compact progress bar")
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) != 1 {
		return errors.New("usage: yt-grab audio <url>")
	}
	url := rest[0]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("url must start with http:// or https://")
	}
	f := choose(*format, cfg.AudioFormat)
	if err := validateAudioFormat(f); err != nil {
		return err
	}
	return runner.Run(runner.Request{URL: url, OutputDir: choose(*output, cfg.OutputDir), AudioOnly: true, AudioFormat: f, PlainOutput: *plain, YtDLPPath: cfg.YtDLPPath, NoPlaylist: true})
}
