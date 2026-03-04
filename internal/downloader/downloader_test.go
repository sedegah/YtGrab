package downloader

import (
	"path/filepath"
	"testing"
)

func TestBuildArgsAudio(t *testing.T) {
	args := BuildArgs(Options{
		URL:       "https://youtube.com/watch?v=abc",
		OutputDir: "/tmp",
		AudioOnly: true,
		AudioFmt:  "mp3",
	})

	wantContains := []string{"-x", "--audio-format", "mp3", "https://youtube.com/watch?v=abc"}
	for _, w := range wantContains {
		if !contains(args, w) {
			t.Fatalf("expected args to contain %q, got %v", w, args)
		}
	}
}

func TestBuildArgsVideoDefault(t *testing.T) {
	args := BuildArgs(Options{
		URL:       "https://youtube.com/watch?v=abc",
		OutputDir: "/tmp",
	})

	if !contains(args, "bestvideo+bestaudio/best") {
		t.Fatalf("expected default best format, got %v", args)
	}
	if !contains(args, filepath.Join("/tmp", "%(title)s.%(ext)s")) {
		t.Fatalf("expected output template path, got %v", args)
	}
}

func contains(s []string, needle string) bool {
	for _, v := range s {
		if v == needle {
			return true
		}
	}
	return false
}
