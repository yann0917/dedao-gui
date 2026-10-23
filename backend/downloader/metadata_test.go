package downloader

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/yann0917/dedao-gui/backend/utils"
)

func TestExistingMP3IsNotDownloadedAgain(t *testing.T) {
	dir := t.TempDir()
	v := Datum{Title: "章节/中文", OrderNum: 7, Type: "audio"}
	name := AudioFilename(v)
	if name != "007.章节 中文.mp3" {
		t.Fatalf("unexpected filename %s", name)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("downloaded audio"), 0644); err != nil {
		t.Fatal(err)
	}
	// No streams or executable: this must return before any network/tool call.
	if err := Download(v, "", dir); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "downloaded audio" {
		t.Fatal("existing MP3 changed")
	}
}

func TestFailedAudioDownloadDoesNotPublishPartialMP3(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable fixture")
	}
	dir := t.TempDir()
	tool := filepath.Join(dir, "ffmpeg")
	if err := os.WriteFile(tool, []byte("#!/bin/sh\nfor last do :; done\nprintf partial > \"$last\"\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	old := utils.FfmpegDir
	utils.FfmpegDir = tool
	t.Cleanup(func() { utils.FfmpegDir = old })
	path := filepath.Join(dir, "track.mp3")
	if err := downloadAudio("https://example.com/audio.m3u8", path); err == nil {
		t.Fatal("expected download failure")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("failed download published a final MP3: %v", err)
	}
	temps, _ := filepath.Glob(filepath.Join(dir, ".audio-*"))
	if len(temps) != 0 {
		t.Fatalf("download leaked temporary files: %v", temps)
	}
}

func TestSuccessfulAudioDownloadIsReadableByMediaServers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable fixture and file modes")
	}
	dir := t.TempDir()
	tool := filepath.Join(dir, "ffmpeg")
	if err := os.WriteFile(tool, []byte("#!/bin/sh\nfor last do :; done\nprintf audio > \"$last\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	old := utils.FfmpegDir
	utils.FfmpegDir = tool
	t.Cleanup(func() { utils.FfmpegDir = old })
	path := filepath.Join(dir, "track.mp3")
	if err := downloadAudio("https://example.com/audio.m3u8", path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0644 {
		t.Fatalf("final MP3 must not retain CreateTemp's private 0600 mode: %v %v", info, err)
	}
}
