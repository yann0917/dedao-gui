package mediameta

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	_ "golang.org/x/image/webp"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

const maxCoverBytes = 10 << 20

// Normalise actual image bytes, not URL extensions (CDNs often return WebP
// even for .jpg URLs). JPEG works in Jellyfin and common ID3 readers alike.
func (e Exporter) cover(ctx context.Context, dir, rawURL string) (string, error) {
	path := filepath.Join(dir, "cover.jpg")
	if data, err := os.ReadFile(path); err == nil {
		if _, format, err := image.DecodeConfig(bytes.NewReader(data)); err == nil && format == "jpeg" {
			return path, nil // preserve a valid local/user-supplied cover on retries
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if rawURL == "" {
		return "", nil
	}
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil {
		return "", fmt.Errorf("无效封面地址")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	client := e.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req) // no account cookies or API credentials
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxCoverBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxCoverBytes {
		return "", fmt.Errorf("封面超过 10 MiB")
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("封面不是有效图片: %w", err)
	}
	if config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 25_000_000 {
		return "", fmt.Errorf("封面尺寸过大或无效")
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	var jpegData bytes.Buffer
	if err = jpeg.Encode(&jpegData, img, &jpeg.Options{Quality: 90}); err != nil {
		return "", err
	}
	if err = atomicWrite(ctx, path, jpegData.Bytes()); err != nil {
		return "", err
	}
	return path, nil
}
