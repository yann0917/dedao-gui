package mediameta

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFFmpegFailurePreservesOriginal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable fixture")
	}
	dir := t.TempDir()
	tool := filepath.Join(dir, "failing-ffmpeg")
	if err := os.WriteFile(tool, []byte("#!/bin/sh\nfor last do :; done\nprintf partial > \"$last\"\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "original.mp3")
	original := []byte("original audio bytes")
	if err := os.WriteFile(path, original, 0640); err != nil {
		t.Fatal(err)
	}
	err := (FFmpegTagger{Executable: tool}).Tag(context.Background(), path, Album{Title: "a"}, Track{Title: "t"}, "")
	if err == nil {
		t.Fatal("expected ffmpeg failure")
	}
	data, _ := os.ReadFile(path)
	if !bytes.Equal(data, original) {
		t.Fatalf("failed tag operation changed original: %s", data)
	}
	temps, _ := filepath.Glob(filepath.Join(dir, ".metadata-*"))
	if len(temps) != 0 {
		t.Fatalf("temp files leaked: %v", temps)
	}
}

func TestFFmpegMetadataRoundTrip(t *testing.T) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("ffmpeg not installed")
	}
	ffprobe, err := exec.LookPath("ffprobe")
	if err != nil {
		t.Skip("ffprobe not installed")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "有空格的 track.mp3")
	if out, err := exec.Command(ffmpeg, "-nostdin", "-y", "-f", "lavfi", "-i", "sine=frequency=440:duration=0.25", "-c:a", "libmp3lame", path).CombinedOutput(); err != nil {
		t.Fatalf("generate audio: %v\n%s", err, out)
	}
	audioHash := func() []byte {
		out, err := exec.Command(ffmpeg, "-v", "error", "-i", path, "-map", "0:a:0", "-c:a", "copy", "-f", "hash", "-hash", "sha256", "-").CombinedOutput()
		if err != nil {
			t.Fatalf("hash audio: %v %s", err, out)
		}
		return out
	}
	before := audioHash()
	cover := filepath.Join(dir, "cover.png")
	if err = os.WriteFile(cover, pngFixture(t), 0644); err != nil {
		t.Fatal(err)
	}
	album := Album{Title: "中文课程 & 标题", Artist: "讲师甲", Description: "课程简介", Date: "2026-10-23"}
	tagger := FFmpegTagger{Executable: ffmpeg}
	for i := 0; i < 2; i++ {
		if err = tagger.Tag(context.Background(), path, album, Track{Title: "第七讲", Number: 7}, cover); err != nil {
			t.Fatal(err)
		}
	}
	// Missing remote artwork must not strip an already embedded local cover.
	if err = tagger.Tag(context.Background(), path, album, Track{Title: "第七讲", Number: 7}, ""); err != nil {
		t.Fatal(err)
	}
	if after := audioHash(); !bytes.Equal(before, after) {
		t.Fatalf("tagging re-encoded or changed audio: before=%s after=%s", before, after)
	}
	out, err := exec.Command(ffprobe, "-v", "error", "-show_format", "-show_streams", "-of", "json", path).Output()
	if err != nil {
		t.Fatal(err)
	}
	var probe struct {
		Format struct {
			Tags map[string]string `json:"tags"`
		} `json:"format"`
		Streams []struct {
			Disposition struct {
				AttachedPic int `json:"attached_pic"`
			} `json:"disposition"`
		} `json:"streams"`
	}
	if err = json.Unmarshal(out, &probe); err != nil {
		t.Fatal(err)
	}
	for key, value := range map[string]string{"title": "第七讲", "album": album.Title, "artist": album.Artist, "album_artist": album.Artist, "track": "7", "date": album.Date, "comment": album.Description} {
		if got := probe.Format.Tags[key]; got != value {
			t.Errorf("%s = %q; want %q", key, got, value)
		}
	}
	attached := 0
	for _, stream := range probe.Streams {
		attached += stream.Disposition.AttachedPic
	}
	if attached != 1 {
		t.Fatalf("retries must replace rather than duplicate embedded covers: %d", attached)
	}
}
