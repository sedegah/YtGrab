package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"yt-grab/internal/config"
)

var version = "dev"

func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	root := flag.NewFlagSet("yt-grab", flag.ContinueOnError)
	configPath := root.String("config", "", "Path to config file (default: ~/.ytgrab.yaml)")
	if err := root.Parse(os.Args[1:]); err != nil {
		return err
	}
	args := root.Args()
	if len(args) == 0 {
		printUsage()
		return nil
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	switch args[0] {
	case "grab":
		return runGrab(cfg, args[1:])
	case "audio":
		return runAudio(cfg, args[1:])
	case "playlist":
		return runPlaylist(cfg, args[1:])
	case "config":
		return runConfig(cfg, args[1:])
	case "doctor":
		return runDoctor()
	case "version":
		fmt.Printf("yt-grab %s\n", version)
		return nil
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return errors.New("unknown command: " + args[0])
	}
}

func printUsage() {
	fmt.Println("yt-grab - A Go-based media downloader and automation CLI")
	fmt.Println("\nUsage:\n  yt-grab [--config path] <command> [flags] <url>")
	fmt.Println("\nCommands:\n  grab      Download a single video\n  audio     Extract audio only\n  playlist  Download playlist\n  config    Manage config (init|view)\n  doctor    Check yt-dlp and ffmpeg\n  version   Print version")
}
