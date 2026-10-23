package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

var FfmpegDir = ""

func runMergeCmd(cmd *exec.Cmd, paths []string, mergeFilePath string) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}

	if mergeFilePath != "" {
		os.Remove(mergeFilePath) // nolint
	}
	// remove parts
	for _, path := range paths {
		os.Remove(path) // nolint
	}
	return nil
}

// MergeAudio merge audio
func MergeAudio(paths []string, mergedFilePath string) error {
	ffmpeg, err := resolveExecutable(FfmpegDir, "ffmpeg")
	if err != nil {
		return err
	}
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", mergedFilePath)
	return runMergeCmd(exec.Command(ffmpeg, cmds...), paths, "")
}

// MergeAudioAndVideo merge audio and video
func MergeAudioAndVideo(paths []string, mergedFilePath string) error {
	ffmpeg, err := resolveExecutable(FfmpegDir, "ffmpeg")
	if err != nil {
		return err
	}
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", "-c:a", "copy", mergedFilePath)
	return runMergeCmd(exec.Command(ffmpeg, cmds...), paths, "")
}

// MergeToMP4 merges video parts to an MP4 file.
func MergeToMP4(paths []string, mergedFilePath string, filename string) error {
	ffmpeg, err := resolveExecutable(FfmpegDir, "ffmpeg")
	if err != nil {
		return err
	}
	mergeFilePath := filename + ".txt" // merge list file should be in the current directory
	// write ffmpeg input file list
	mergeFile, err := os.Create(mergeFilePath)
	if err != nil {
		return err
	}
	for _, path := range paths {
		_, _ = mergeFile.Write([]byte(fmt.Sprintf("file '%s'\n", path))) // nolint
	}
	err = mergeFile.Close() // nolint
	if err != nil {
		return err
	}

	cmd := exec.Command(
		ffmpeg, "-y", "-f", "concat", "-safe", "-1",
		"-i", mergeFilePath, "-c", "copy", "-bsf:a", "aac_adtstoasc", mergedFilePath,
	)
	return runMergeCmd(cmd, paths, mergeFilePath)
}
