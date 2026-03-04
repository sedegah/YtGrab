package cli

import (
	"fmt"
	"runtime"

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
		printDependencyHints()
		return fmt.Errorf("dependency check failed")
	}
	return nil
}

func printDependencyHints() {
	fmt.Println("\nInstall missing dependencies:")
	switch runtime.GOOS {
	case "windows":
		fmt.Println("  winget install yt-dlp.yt-dlp")
		fmt.Println("  winget install Gyan.FFmpeg")
		fmt.Println("  Restart PowerShell after installation.")
	case "darwin":
		fmt.Println("  brew install yt-dlp ffmpeg")
	default:
		fmt.Println("  sudo apt update && sudo apt install -y yt-dlp ffmpeg")
	}
}
