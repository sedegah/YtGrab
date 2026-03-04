package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	OutputDir     string
	Format        string
	AudioFormat   string
	MaxResolution int
	YtDLPPath     string
}

func Default() Config {
	home, _ := os.UserHomeDir()
	return Config{
		OutputDir:     filepath.Join(home, "Downloads"),
		Format:        "bestvideo+bestaudio/best",
		AudioFormat:   "mp3",
		MaxResolution: 0,
		YtDLPPath:     "",
	}
}

func Load(configFile string) (Config, error) {
	cfg := Default()
	path, err := resolvePath(configFile)
	if err != nil {
		return cfg, err
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		switch k {
		case "output_dir":
			cfg.OutputDir = expandHome(v)
		case "format":
			cfg.Format = v
		case "audio_format":
			cfg.AudioFormat = v
		case "max_resolution":
			n, _ := strconv.Atoi(v)
			cfg.MaxResolution = n
		case "yt_dlp_path":
			cfg.YtDLPPath = v
		}
	}
	if err := s.Err(); err != nil {
		return cfg, err
	}
	if env := os.Getenv("YTGRAB_OUTPUT_DIR"); env != "" {
		cfg.OutputDir = env
	}
	if env := os.Getenv("YTGRAB_FORMAT"); env != "" {
		cfg.Format = env
	}
	if env := os.Getenv("YTGRAB_AUDIO_FORMAT"); env != "" {
		cfg.AudioFormat = env
	}
	return cfg, nil
}

func ConfigPath() (string, error) {
	return resolvePath("")
}

func resolvePath(configFile string) (string, error) {
	if configFile != "" {
		return configFile, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	return filepath.Join(home, ".ytgrab.yaml"), nil
}

func expandHome(v string) string {
	if strings.HasPrefix(v, "~/") {
		h, _ := os.UserHomeDir()
		return filepath.Join(h, strings.TrimPrefix(v, "~/"))
	}
	return v
}

func (c Config) ToYAML() string {
	return fmt.Sprintf("output_dir: %s\nformat: %s\naudio_format: %s\nmax_resolution: %d\nyt_dlp_path: %s\n", c.OutputDir, c.Format, c.AudioFormat, c.MaxResolution, c.YtDLPPath)
}
