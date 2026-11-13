package downloader

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/yann0917/dedao-gui/backend/utils"
)

func TestUnverifiableExistingMP3IsNotClobbered(t *testing.T) {
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
	// 没有时长信息无法校验完整性，也没有播放地址可重新下载：
	// 必须报错返回，且不能动已有文件。
	if _, err := Download(v, "", dir); err == nil {
		t.Fatal("expected error for unverifiable file without m3u8 url")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "downloaded audio" {
		t.Fatal("existing MP3 changed")
	}
}

func TestVerifiedExistingMP3IsSkipped(t *testing.T) {
	if !audioProbeAvailable {
		t.Skip("本机无 ffprobe，跳过")
	}
	dir := t.TempDir()
	v := Datum{Title: "章节/中文", OrderNum: 7, Type: "audio", Duration: 2}
	path := filepath.Join(dir, AudioFilename(v))
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "lavfi",
		"-i", "sine=frequency=440:duration=2",
		"-q:a", "9",
		path,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("生成测试音频失败: %v\n%s", err, out)
	}

	skipped, err := Download(v, "", dir)
	if err != nil {
		t.Fatalf("完整音频应被跳过: %v", err)
	}
	if !skipped {
		t.Fatal("完整音频应返回 skipped=true")
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
	v := Datum{Title: "track", Type: "audio", M3U8URL: "https://example.com/audio.m3u8"}
	if _, err := downloadAudio(v, path); err == nil {
		t.Fatal("expected download failure")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("failed download published a final MP3: %v", err)
	}
	temps, _ := filepath.Glob(filepath.Join(dir, ".dedao-*.part.mp3"))
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
	v := Datum{Title: "track", Type: "audio", M3U8URL: "https://example.com/audio.m3u8"}
	if _, err := downloadAudio(v, path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0644 {
		t.Fatalf("final MP3 must not retain CreateTemp's private 0600 mode: %v %v", info, err)
	}
}
