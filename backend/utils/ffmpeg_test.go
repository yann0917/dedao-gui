package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMergeAudioFindsFFmpegWhenSettingEmpty(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses a POSIX executable fixture")
	}
	dir := t.TempDir()
	marker := filepath.Join(dir, "called")
	tool := filepath.Join(dir, "ffmpeg")
	if err := os.WriteFile(tool, []byte("#!/bin/sh\nprintf called > \"$FFMPEG_TEST_MARKER\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("FFMPEG_TEST_MARKER", marker)
	old := FfmpegDir
	FfmpegDir = ""
	t.Cleanup(func() { FfmpegDir = old })
	if err := MergeAudio([]string{filepath.Join(dir, "part.ts")}, filepath.Join(dir, "out.mp3")); err != nil {
		t.Fatalf("empty setting should discover ffmpeg instead of exec: no command: %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("ffmpeg was not invoked: %v", err)
	}
}
