package mediameta

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type tagCall struct {
	path  string
	album Album
	track Track
	cover string
}

type recordingTagger struct {
	calls []tagCall
	err   error
}

func (t *recordingTagger) Tag(_ context.Context, path string, album Album, track Track, cover string) error {
	t.calls = append(t.calls, tagCall{path, album, track, cover})
	return t.err
}

func pngFixture(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	img.Set(0, 0, color.White)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func readNFO(t *testing.T, dir string) albumNFO {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "album.nfo"))
	if err != nil {
		t.Fatal(err)
	}
	var nfo albumNFO
	if err = xml.Unmarshal(data, &nfo); err != nil {
		t.Fatalf("invalid XML: %v\n%s", err, data)
	}
	return nfo
}

func TestExportSidecarsAndExistingTracks(t *testing.T) {
	img := pngFixture(t)
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Header.Get("Cookie") != "" || r.Header.Get("Authorization") != "" {
			t.Error("cover request leaked account credentials")
		}
		w.Write(img) // intentionally no extension in URL, PNG rather than JPEG
	}))
	defer server.Close()
	dir := t.TempDir()
	name := "007.章节 <问答>.mp3"
	if err := os.WriteFile(filepath.Join(dir, name), []byte("audio"), 0644); err != nil {
		t.Fatal(err)
	}
	tagger := &recordingTagger{}
	e := Exporter{HTTP: server.Client(), Tagger: tagger}
	album := Album{Title: "课程 & <文明>", Artist: "甲 & 乙", Description: "中文简介\n第二行", Date: "2026-10-23", ID: "123", CoverURL: server.URL + "/image?token=secret"}
	tracks := []Track{{Filename: name, Title: "章节 <问答>", Number: 7}, {Filename: "missing.mp3"}}
	for i := 0; i < 2; i++ {
		if err := e.Export(context.Background(), dir, album, tracks); err != nil {
			t.Fatal(err)
		}
	}
	nfo := readNFO(t, dir)
	if nfo.Title != album.Title || nfo.Artist != album.Artist || nfo.Review != album.Description || nfo.Date != album.Date || nfo.ID == nil || nfo.ID.Type != "dedao" || nfo.ID.Value != album.ID {
		t.Fatalf("unexpected NFO: %+v", nfo)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "album.nfo"))
	if strings.Contains(string(data), "secret") || strings.Contains(string(data), "http:") || strings.Contains(string(data), "<thumb") || strings.Contains(string(data), dir) {
		t.Fatalf("NFO must be portable, not contain signed source URLs: %s", data)
	}
	if requests != 1 || len(tagger.calls) != 2 {
		t.Fatalf("retry should reuse cover and only tag existing tracks: requests=%d calls=%d", requests, len(tagger.calls))
	}
	call := tagger.calls[0]
	if call.track != tracks[0] || call.album != album || call.cover != filepath.Join(dir, "cover.jpg") {
		t.Fatalf("wrong tag mapping: %+v", call)
	}
	cover, _ := os.ReadFile(call.cover)
	if _, format, err := image.Decode(bytes.NewReader(cover)); err != nil || format != "jpeg" {
		t.Fatalf("cover.jpg is not a real JPEG: %s %v", format, err)
	}
}

func TestMetadataOnlyWithoutAudioNeedsNoFFmpegOrCover(t *testing.T) {
	dir := t.TempDir()
	if err := (Exporter{}).Export(context.Background(), dir, Album{Title: "标题"}, []Track{{Filename: "missing.mp3"}}); err != nil {
		t.Fatal(err)
	}
	nfo := readNFO(t, dir)
	if nfo.Title != "标题" || nfo.Date != "" || nfo.ID != nil {
		t.Fatalf("unknown metadata must not be fabricated: %+v", nfo)
	}
}

func TestInvalidCoverDoesNotPublishPartialFiles(t *testing.T) {
	for _, tc := range []struct {
		name string
		code int
		body string
	}{
		{"http-error", 403, "forbidden"},
		{"html-not-image", 200, "<html>login required</html>"},
		{"truncated-image", 200, "\x89PNG\r\n\x1a\n"},
		{"oversize", 200, strings.Repeat("x", maxCoverBytes+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.code)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()
			dir := t.TempDir()
			err := (Exporter{}).Export(context.Background(), dir, Album{Title: "a", CoverURL: server.URL}, nil)
			if err == nil {
				t.Fatal("expected cover error")
			}
			entries, _ := os.ReadDir(dir)
			if len(entries) != 0 {
				t.Fatalf("failed cover must not leave final or temp files: %v", entries)
			}
		})
	}
}

func TestExportCancellation(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := (Exporter{}).Export(ctx, dir, Album{}, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancel()
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	if err := (Exporter{}).Export(ctx, dir, Album{CoverURL: server.URL}, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("HTTP request did not respect cancellation: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatal("canceled export left files")
	}
}

func TestTaggingFailureIsVisibleAndPathsStayInsideAlbum(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "track.mp3"), []byte("audio"), 0644); err != nil {
		t.Fatal(err)
	}
	tagErr := errors.New("tag failed")
	e := Exporter{Tagger: &recordingTagger{err: tagErr}}
	err := e.Export(context.Background(), dir, Album{}, []Track{{Filename: "track.mp3"}})
	if !errors.Is(err, tagErr) {
		t.Fatalf("tagger failure was swallowed: %v", err)
	}
	for _, filename := range []string{"../outside.mp3", "/outside.mp3", ""} {
		if err = e.Export(context.Background(), dir, Album{}, []Track{{Filename: filename}}); err == nil {
			t.Fatalf("accepted unsafe filename %q", filename)
		}
	}
}
