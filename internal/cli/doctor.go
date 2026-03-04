package cli

import (
	"fmt"

	"yt-grab/internal/runner"
)

func runDoctor() error {
	deps := []string{"yt-dlp", "ffmpeg"}
	failed := false
	for _, d := range deps {
		if err := runner.CheckDependency(d); err != nil {
			fmt.Printf("✗ %s\n", err)
			failed = true
		} else {
			fmt.Printf("✓ %s found\n", d)
		}
	}
	if failed {
		return fmt.Errorf("dependency check failed")
	}
	return nil
}
