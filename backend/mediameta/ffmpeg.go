package mediameta

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type FFmpegTagger struct {
	Executable string
}

// Tag remuxes without re-encoding the audio. Only replace the original after a
// successful command and cancellation check; a failure leaves it intact.
func (t FFmpegTagger) Tag(ctx context.Context, path string, album Album, track Track, cover string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tool := t.Executable
	if tool == "" {
		tool = "ffmpeg"
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".metadata-*.mp3")
	if err != nil {
		return err
	}
	temp := f.Name()
	defer os.Remove(temp)
	if err = f.Close(); err != nil {
		return err
	}
	args := []string{"-nostdin", "-y", "-i", path}
	if cover != "" {
		args = append(args, "-i", cover)
	}
	args = append(args, "-map", "0:a:0")
	if cover != "" {
		args = append(args, "-map", "1:v:0", "-c:v", "copy", "-disposition:v:0", "attached_pic")
	} else {
		args = append(args, "-map", "0:v?", "-c:v", "copy")
	}
	args = append(args, "-map_metadata", "0", "-c:a", "copy", "-id3v2_version", "3")
	fields := [][2]string{
		{"title", track.Title}, {"album", album.Title},
		{"artist", album.Artist}, {"album_artist", album.Artist},
		{"comment", album.Description}, {"date", album.Date},
	}
	if track.Number > 0 {
		fields = append(fields, [2]string{"track", strconv.Itoa(track.Number)})
	}
	for _, field := range fields {
		if field[1] != "" {
			args = append(args, "-metadata", field[0]+"="+field[1])
		}
	}
	args = append(args, temp)
	cmd := exec.CommandContext(ctx, tool, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("ffmpeg: %w\n%s", err, stderr.String())
	}
	info, err := os.Stat(temp)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("ffmpeg 未生成有效的音频文件")
	}
	original, err := os.Stat(path)
	if err != nil {
		return err
	}
	if err = os.Chmod(temp, original.Mode().Perm()); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return os.Rename(temp, path)
}
