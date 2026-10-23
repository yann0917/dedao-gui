// Package mediameta exports portable local audio metadata without depending on
// the Dedao API, desktop runtime, or download queue.
package mediameta

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type Album struct {
	Title       string
	Artist      string
	Description string
	Date        string // YYYY-MM-DD, absent when the source has no publication date
	ID          string
	CoverURL    string
}

type Track struct {
	Filename string // basename only, shared with the audio downloader
	Title    string
	Number   int
}

// AudioTagger is implemented by the ffmpeg adapter; tests and other containers
// can supply their own implementation at the composition root.
type AudioTagger interface {
	Tag(context.Context, string, Album, Track, string) error
}

type Exporter struct {
	HTTP   HTTPClient
	Tagger AudioTagger
}

type albumNFO struct {
	XMLName xml.Name    `xml:"album"`
	Title   string      `xml:"title"`
	Artist  string      `xml:"artist,omitempty"`
	Review  string      `xml:"review,omitempty"`
	Date    string      `xml:"releasedate,omitempty"`
	ID      *providerID `xml:"uniqueid,omitempty"`
}

type providerID struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// Export always writes sidecars, but only tags audio files that already exist.
// This also supports metadata-only tasks without downloading audio streams.
func (e Exporter) Export(ctx context.Context, dir string, album Album, tracks []Track) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	cover, err := e.cover(ctx, dir, album.CoverURL)
	if err != nil {
		return fmt.Errorf("下载封面: %w", err)
	}
	nfo := albumNFO{Title: album.Title, Artist: album.Artist, Review: album.Description, Date: album.Date}
	if album.ID != "" {
		nfo.ID = &providerID{Type: "dedao", Value: album.ID}
	}
	// Jellyfin's NFO image parser requires absolute paths/URLs. Neither an
	// exporter-machine path nor a signed CDN URL is portable to a media server;
	// let the standard local cover.jpg discovery supply the image instead.
	data, err := xml.MarshalIndent(nfo, "", "  ")
	if err != nil {
		return err
	}
	if err = atomicWrite(ctx, filepath.Join(dir, "album.nfo"), append([]byte(xml.Header), data...)); err != nil {
		return fmt.Errorf("写入 NFO: %w", err)
	}
	for _, track := range tracks {
		if err := ctx.Err(); err != nil {
			return err
		}
		if track.Filename == "" || filepath.Base(track.Filename) != track.Filename {
			return fmt.Errorf("无效音频文件名: %q", track.Filename)
		}
		path := filepath.Join(dir, track.Filename)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Size() == 0 {
			return fmt.Errorf("无效音频文件: %s", path)
		}
		if e.Tagger == nil {
			return fmt.Errorf("未配置音频标签写入器")
		}
		if err := e.Tagger.Tag(ctx, path, album, track, cover); err != nil {
			return fmt.Errorf("写入音频元数据 %s: %w", track.Filename, err)
		}
	}
	return nil
}

func atomicWrite(ctx context.Context, path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".metadata-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Chmod(0644); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
