package cli

import (
	"errors"
	"fmt"
	"os"

	"yt-grab/internal/config"
)

func runConfig(cfg config.Config, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: yt-grab config <init|view>")
	}
	switch args[0] {
	case "init":
		path, err := config.ConfigPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config already exists at %s", path)
		}
		return os.WriteFile(path, []byte(config.Default().ToYAML()), 0o644)
	case "view":
		fmt.Print(cfg.ToYAML())
		return nil
	default:
		return errors.New("usage: yt-grab config <init|view>")
	}
}
