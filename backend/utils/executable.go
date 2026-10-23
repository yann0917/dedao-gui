package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Explicit settings override PATH; Finder does not normally inherit shell PATH.
func resolveExecutable(configured, name string) (string, error) {
	var fallbackDirs []string
	if runtime.GOOS == "darwin" {
		fallbackDirs = []string{"/opt/homebrew/bin", "/usr/local/bin"}
	}
	return findExecutable(configured, name, fallbackDirs)
}

func findExecutable(configured, name string, fallbackDirs []string) (string, error) {
	if strings.TrimSpace(configured) != "" {
		path, err := exec.LookPath(configured)
		if err != nil {
			return "", fmt.Errorf("%s 路径无效或不可执行，请在系统设置中检查路径 %q: %w", name, configured, err)
		}
		return path, nil
	}
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	for _, dir := range fallbackDirs {
		if path, err := exec.LookPath(filepath.Join(dir, name)); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到 %s，请先安装，并在系统设置中填写其可执行文件完整路径", name)
}
