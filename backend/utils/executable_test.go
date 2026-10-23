package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFindExecutable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses POSIX executable fixtures")
	}
	pathDir, fallbackDir, explicitDir := t.TempDir(), t.TempDir(), t.TempDir()
	for _, dir := range []string{pathDir, fallbackDir, explicitDir} {
		if err := os.WriteFile(filepath.Join(dir, "ffmpeg"), []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", pathDir)
	for _, tc := range []struct {
		name, configured, path, want string
	}{
		{"explicit override", filepath.Join(explicitDir, "ffmpeg"), pathDir, filepath.Join(explicitDir, "ffmpeg")},
		{"PATH before fallback", "", pathDir, filepath.Join(pathDir, "ffmpeg")},
		{"Finder PATH fallback", "", t.TempDir(), filepath.Join(fallbackDir, "ffmpeg")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PATH", tc.path)
			got, err := findExecutable(tc.configured, "ffmpeg", []string{fallbackDir})
			if err != nil || got != tc.want {
				t.Fatalf("got %q, %v; want %q", got, err, tc.want)
			}
		})
	}
	t.Run("invalid override must not silently fall back", func(t *testing.T) {
		_, err := findExecutable(filepath.Join(explicitDir, "missing"), "ffmpeg", []string{fallbackDir})
		if err == nil || !strings.Contains(err.Error(), "系统设置") {
			t.Fatalf("expected actionable error, got %v", err)
		}
	})
	t.Run("missing tool", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		_, err := findExecutable("", "ffmpeg", nil)
		if err == nil || !strings.Contains(err.Error(), "ffmpeg") || !strings.Contains(err.Error(), "系统设置") {
			t.Fatalf("expected actionable error, got %v", err)
		}
	})
}
