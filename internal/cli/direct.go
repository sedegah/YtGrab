package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"yt-grab/internal/config"
	"yt-grab/internal/runner"
)

func runDirect(cfg config.Config, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: yt-grab <url> [--audio] [--quality best|worst|720p|1080p]")
	}

	url := args[0]
	if !isURL(url) {
		return errors.New("url must start with http:// or https://")
	}

	output := cfg.OutputDir
	format := cfg.Format
	audioOnly := false
	audioFormat := cfg.AudioFormat
	maxRes := cfg.MaxResolution
	quality := ""
	formatSet := false

	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--audio":
			audioOnly = true
		case "--output", "-o":
			if i+1 >= len(args) {
				return fmt.Errorf("missing value for %s", arg)
			}
			i++
			output = args[i]
		case "--audio-format":
			if i+1 >= len(args) {
				return errors.New("missing value for --audio-format")
			}
			i++
			audioFormat = args[i]
		case "--format":
			if i+1 >= len(args) {
				return errors.New("missing value for --format")
			}
			i++
			format = args[i]
			formatSet = true
		case "--max-res":
			if i+1 >= len(args) {
				return errors.New("missing value for --max-res")
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil || n <= 0 {
				return errors.New("--max-res must be a positive number")
			}
			maxRes = n
		case "--quality", "-q":
			if i+1 >= len(args) {
				return errors.New("missing value for --quality")
			}
			i++
			quality = args[i]
		case "-h", "--help", "help":
			printUsage()
			return nil
		default:
			return fmt.Errorf("unknown flag: %s", arg)
		}
	}

	if err := validateAudioFormat(audioFormat); err != nil {
		return err
	}

	if quality != "" {
		if formatSet {
			return errors.New("use either --format or --quality, not both")
		}
		qFormat, qMaxRes, err := qualityToFormat(quality)
		if err != nil {
			return err
		}
		if qFormat != "" {
			format = qFormat
			maxRes = 0
		}
		if qMaxRes > 0 {
			maxRes = qMaxRes
		}
	}

	return runner.Run(runner.Request{
		URL:           url,
		OutputDir:     output,
		Format:        format,
		AudioOnly:     audioOnly,
		AudioFormat:   audioFormat,
		MaxResolution: maxRes,
		YtDLPPath:     cfg.YtDLPPath,
		NoPlaylist:    true,
	})
}

func validateAudioFormat(v string) error {
	if v != "mp3" && v != "aac" && v != "wav" {
		return errors.New("audio format must be one of: mp3, aac, wav")
	}
	return nil
}

func qualityToFormat(v string) (string, int, error) {
	q := strings.ToLower(strings.TrimSpace(v))
	switch q {
	case "best":
		return "bestvideo+bestaudio/best", 0, nil
	case "worst":
		return "worstvideo+worstaudio/worst", 0, nil
	}
	if strings.HasSuffix(q, "p") {
		n, err := strconv.Atoi(strings.TrimSuffix(q, "p"))
		if err == nil && n > 0 {
			return "", n, nil
		}
	}
	return "", 0, errors.New("--quality must be best, worst, or a resolution like 720p or 1080p")
}
