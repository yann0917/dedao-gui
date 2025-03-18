package utils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var FfmpegDir = ""

// ProgressCallback 是一个回调函数类型，用于报告ffmpeg处理进度
type ProgressCallback func(percentage int, status string)

// runMergeCmd 执行ffmpeg命令，支持进度回调
func runMergeCmd(ctx context.Context, cmd *exec.Cmd, paths []string, mergeFilePath string, callback ProgressCallback) error {
	var stderr bytes.Buffer
	
	// 在Windows平台隐藏控制台窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	
	// 创建stderr管道用于捕获输出
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("create stderr pipe: %v", err)
	}
	
	// 同时将输出存储到buffer以便错误处理
	cmd.Stderr = io.MultiWriter(&stderr, stderrPipe)
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %v", err)
	}
	
	// 正则表达式匹配ffmpeg输出中的时间信息
	durationRegex := regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2})`)
	timeRegex := regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2})`)
	
	// 处理进度输出
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		var totalSeconds float64
		
		for scanner.Scan() {
			line := scanner.Text()
			
			// 解析总时长
			if totalSeconds == 0 {
				if matches := durationRegex.FindStringSubmatch(line); len(matches) > 3 {
					h, _ := strconv.ParseFloat(matches[1], 64)
					m, _ := strconv.ParseFloat(matches[2], 64)
					s, _ := strconv.ParseFloat(matches[3], 64)
					totalSeconds = h*3600 + m*60 + s
				}
			}
			
			// 解析当前处理时间
			if matches := timeRegex.FindStringSubmatch(line); len(matches) > 3 && totalSeconds > 0 {
				h, _ := strconv.ParseFloat(matches[1], 64)
				m, _ := strconv.ParseFloat(matches[2], 64)
				s, _ := strconv.ParseFloat(matches[3], 64)
				currentSeconds := h*3600 + m*60 + s
				
				// 计算百分比进度
				percentage := int((currentSeconds / totalSeconds) * 100)
				if percentage > 100 {
					percentage = 100
				}
				
				// 调用回调函数报告进度
				if callback != nil {
					callback(percentage, fmt.Sprintf("处理中 %.0f%%", float64(percentage)))
				}
			}
		}
	}()
	
	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}
	
	// 执行完成后报告100%进度
	if callback != nil {
		callback(100, "处理完成")
	}
	
	// 清理临时文件
	if mergeFilePath != "" {
		os.Remove(mergeFilePath) // nolint
	}
	// 删除分片文件
	for _, path := range paths {
		os.Remove(path) // nolint
	}
	return nil
}

// MergeAudio 合并音频文件，支持进度回调
func MergeAudio(ctx context.Context, paths []string, mergedFilePath string, callback ProgressCallback) error {
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", mergedFilePath)
	return runMergeCmd(ctx, exec.Command(FfmpegDir, cmds...), paths, "", callback)
}

// MergeAudioAndVideo 合并音频和视频文件，支持进度回调
func MergeAudioAndVideo(ctx context.Context, paths []string, mergedFilePath string, callback ProgressCallback) error {
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", "-c:a", "copy", mergedFilePath)
	return runMergeCmd(ctx, exec.Command(FfmpegDir, cmds...), paths, "", callback)
}

// MergeToMP4 将视频片段合并为MP4文件，支持进度回调
func MergeToMP4(ctx context.Context, paths []string, mergedFilePath string, filename string, callback ProgressCallback) error {
	mergeFilePath := filename + ".txt" // merge list file should be in the current directory
	// write ffmpeg input file list
	mergeFile, _ := os.Create(mergeFilePath)
	for _, path := range paths {
		_, _ = mergeFile.Write([]byte(fmt.Sprintf("file '%s'\n", path))) // nolint
	}
	err := mergeFile.Close() // nolint
	if err != nil {
		return err
	}

	cmd := exec.Command(
		FfmpegDir, "-y", "-f", "concat", "-safe", "-1",
		"-i", mergeFilePath, "-c", "copy", "-bsf:a", "aac_adtstoasc", mergedFilePath,
	)
	return runMergeCmd(ctx, cmd, paths, mergeFilePath, callback)
}

// 为了向后兼容，提供不带进度回调的旧API
func MergeAudioLegacy(paths []string, mergedFilePath string) error {
	return MergeAudio(context.Background(), paths, mergedFilePath, nil)
}

func MergeAudioAndVideoLegacy(paths []string, mergedFilePath string) error {
	return MergeAudioAndVideo(context.Background(), paths, mergedFilePath, nil)
}

func MergeToMP4Legacy(paths []string, mergedFilePath string, filename string) error {
	return MergeToMP4(context.Background(), paths, mergedFilePath, filename, nil)
}

// 向后兼容的API覆盖
func MergeAudio(paths []string, mergedFilePath string) error {
	return MergeAudioLegacy(paths, mergedFilePath)
}

func MergeAudioAndVideo(paths []string, mergedFilePath string) error {
	return MergeAudioAndVideoLegacy(paths, mergedFilePath)
}

func MergeToMP4(paths []string, mergedFilePath string, filename string) error {
	return MergeToMP4Legacy(paths, mergedFilePath, filename)
}
