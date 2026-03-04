package runner

import "testing"

func TestBuildArgsAudio(t *testing.T) {
	args := BuildArgs(Request{
		URL:         "https://youtube.com/watch?v=abc",
		OutputDir:   "/tmp",
		AudioOnly:   true,
		AudioFormat: "mp3",
	})
	mustContain(t, args, "-x")
	mustContain(t, args, "--audio-format")
	mustContain(t, args, "mp3")
}

func TestBuildArgsMaxResolution(t *testing.T) {
	args := BuildArgs(Request{
		URL:           "https://youtube.com/watch?v=abc",
		OutputDir:     "/tmp",
		MaxResolution: 1080,
	})
	mustContain(t, args, "bestvideo[height<=1080]+bestaudio/best[height<=1080]")
}

func mustContain(t *testing.T, args []string, expected string) {
	t.Helper()
	for _, arg := range args {
		if arg == expected {
			return
		}
	}
	t.Fatalf("expected %q in %v", expected, args)
}
