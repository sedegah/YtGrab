package runner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Request struct {
	URL           string
	OutputDir     string
	Format        string
	AudioOnly     bool
	AudioFormat   string
	PlainOutput   bool
	MaxResolution int
	NoPlaylist    bool
	FlatPlaylist  bool
	YtDLPPath     string
}

func BuildArgs(req Request) []string {
	args := []string{"--no-part"}
	if !req.PlainOutput {
		args = append(args,
			"--newline",
			"--progress",
			"--progress-template",
			"download:%(progress._percent_str)s|%(progress._speed_str)s|%(progress._eta_str)s",
		)
	}
	args = append(args, "-o", filepath.Join(req.OutputDir, "%(title)s.%(ext)s"))

	if req.NoPlaylist {
		args = append(args, "--no-playlist")
	}
	if req.FlatPlaylist {
		args = append(args, "--flat-playlist")
	}

	if req.AudioOnly {
		args = append(args, "-x", "--audio-format", req.AudioFormat)
	} else {
		format := req.Format
		if req.MaxResolution > 0 {
			cap := strconv.Itoa(req.MaxResolution)
			format = "bestvideo[height<=" + cap + "]+bestaudio/best[height<=" + cap + "]"
		}
		if format == "" {
			format = "bestvideo+bestaudio/best"
		}
		args = append(args, "-f", format)
	}

	return append(args, req.URL)
}

func Run(req Request) error {
	bin := req.YtDLPPath
	if bin == "" {
		resolved, err := exec.LookPath("yt-dlp")
		if err != nil {
			return fmt.Errorf("yt-dlp not found in PATH: %w", err)
		}
		bin = resolved
	}

	cmd := exec.Command(bin, BuildArgs(req)...)
	if req.PlainOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	renderer := &progressRenderer{}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		streamOutput(stdout, renderer)
	}()
	go func() {
		defer wg.Done()
		streamOutput(stderr, renderer)
	}()
	wg.Wait()
	renderer.finish()

	return cmd.Wait()
}

func CheckDependency(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s not found in PATH", name)
	}
	return nil
}

var percentRe = regexp.MustCompile(`(\d+(?:\.\d+)?)%`)

type progressRenderer struct {
	mu             sync.Mutex
	progressActive bool
}

func (p *progressRenderer) printLine(line string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.progressActive {
		fmt.Fprint(os.Stdout, "\n")
		p.progressActive = false
	}
	fmt.Fprintln(os.Stdout, line)
}

func (p *progressRenderer) update(percent float64, speed, eta string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	barWidth := 24
	filled := int((percent / 100) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)
	if speed == "" {
		speed = "-"
	}
	if eta == "" {
		eta = "-"
	}
	fmt.Fprintf(os.Stdout, "\r[%s] %5.1f%%  %s  ETA %s", bar, percent, speed, eta)
	p.progressActive = true
}

func (p *progressRenderer) finish() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.progressActive {
		fmt.Fprint(os.Stdout, "\n")
		p.progressActive = false
	}
}

func streamOutput(r io.Reader, renderer *progressRenderer) {
	s := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	s.Buffer(buf, 1024*1024)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		if pct, speed, eta, ok := parseProgress(line); ok {
			renderer.update(pct, speed, eta)
			continue
		}

		if strings.HasPrefix(line, "[download]") && strings.Contains(line, "%") {
			continue
		}

		renderer.printLine(line)
	}
}

func parseProgress(line string) (float64, string, string, bool) {
	if strings.HasPrefix(line, "download:") {
		parts := strings.SplitN(strings.TrimPrefix(line, "download:"), "|", 3)
		if len(parts) == 3 {
			pct, ok := parsePercent(parts[0])
			if !ok {
				return 0, "", "", false
			}
			return pct, strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), true
		}
	}

	if strings.HasPrefix(line, "[download]") && strings.Contains(line, "%") {
		pct, ok := parsePercent(line)
		if !ok {
			return 0, "", "", false
		}
		speed := ""
		eta := ""
		if idx := strings.Index(line, " at "); idx >= 0 {
			rest := line[idx+4:]
			if etaIdx := strings.Index(rest, " ETA "); etaIdx >= 0 {
				speed = strings.TrimSpace(rest[:etaIdx])
				eta = strings.TrimSpace(rest[etaIdx+5:])
			} else {
				speed = strings.TrimSpace(rest)
			}
		}
		return pct, speed, eta, true
	}

	return 0, "", "", false
}

func parsePercent(v string) (float64, bool) {
	m := percentRe.FindStringSubmatch(v)
	if len(m) < 2 {
		return 0, false
	}
	n, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
