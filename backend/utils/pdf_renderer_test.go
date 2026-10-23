package utils

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

type fakePDFRenderer struct{ err error }

func (r fakePDFRenderer) Render(_ context.Context, _ *PdfOption, _ []byte) error { return r.err }

func TestPDFRendererSelection(t *testing.T) {
	for _, engine := range []string{"", "chromium", "wkhtmltopdf", "invalid"} {
		t.Run(engine, func(t *testing.T) {
			t.Setenv("DEDAO_PDF_ENGINE", engine)
			r, err := selectPDFRenderer()
			if engine == "invalid" {
				if err == nil {
					t.Fatal("unknown engine accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			_, legacy := r.(legacyPDFRenderer)
			if legacy != (engine == "wkhtmltopdf") {
				t.Fatal("unexpected automatic fallback")
			}
		})
	}
}

func TestPDFRendererFailureKeepsCover(t *testing.T) {
	cover := filepath.Join(t.TempDir(), "cover.html")
	if err := os.WriteFile(cover, []byte("cover"), 0600); err != nil {
		t.Fatal(err)
	}
	want := errors.New("render failure")
	p := PdfOption{CoverPath: cover, renderer: fakePDFRenderer{err: want}}
	if err := p.GenPdf(bytes.NewBufferString("body")); !errors.Is(err, want) {
		t.Fatal(err)
	}
	if _, err := os.Stat(cover); err != nil {
		t.Fatal("lost retry input", err)
	}
}

func TestChromiumDiscovery(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := locateChromium("", nil)
	if err == nil || !strings.Contains(err.Error(), "系统设置") {
		t.Fatal("missing actionable error", err)
	}
	_, err = locateChromium(filepath.Join(t.TempDir(), "missing"), []string{"/bin/sh"})
	if err == nil {
		t.Fatal("invalid explicit path silently fell back")
	}
}

func TestChromiumDocument(t *testing.T) {
	doc, err := chromiumDocument([]byte(`<html><head><style>p{color:red}</style></head><body><h1 id="chapter">中文标题</h1><a href="#chapter">目录</a><script>throw new Error('must not run')</script><meta http-equiv="refresh" content="0;url=https://example.com"><img onerror="alert(1)"></body></html>`), []byte(`<html><head><style>body{margin:999mm}</style></head><body><img src="data:image/png;base64,AA=="></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	s := string(doc)
	for _, want := range []string{"script-src 'none'", "@page dedao-cover", "counter(page)", "中文标题", `href="#chapter"`, "p{color:red}"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(s, "margin:999mm") {
		t.Fatal("cover styles leaked into main document")
	}
	if strings.Contains(s, "<script>") {
		t.Fatal("active script retained")
	}
	if strings.Contains(s, "onerror") || strings.Contains(s, `http-equiv="refresh"`) {
		t.Fatal("active HTML retained")
	}
}

func TestChromiumFailedOutputAndCancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses a POSIX fake executable")
	}
	old := ChromiumPath
	t.Cleanup(func() { ChromiumPath = old })
	for _, tc := range []struct {
		name, script string
		cancel       bool
	}{
		{"missing output", "exit 0", false},
		{"browser failure", "echo browser-failed >&2; exit 1", false},
		{"cancel process tree", "sleep 30", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			ChromiumPath = filepath.Join(dir, "browser")
			if err := os.WriteFile(ChromiumPath, []byte("#!/bin/sh\n"+tc.script+"\n"), 0700); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(dir, "existing.pdf")
			if err := os.WriteFile(out, []byte("keep-existing"), 0600); err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			if tc.cancel {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 100*time.Millisecond)
				defer cancel()
			}
			start := time.Now()
			err := (chromiumPDFRenderer{}).Render(ctx, &PdfOption{FileName: out}, []byte("<p>test</p>"))
			if err == nil {
				t.Fatal("browser failure accepted")
			}
			if tc.cancel && time.Since(start) > 5*time.Second {
				t.Fatal("browser process tree did not terminate promptly")
			}
			data, _ := os.ReadFile(out)
			if string(data) != "keep-existing" {
				t.Fatal("failed render overwrote original PDF")
			}
			staged, _ := filepath.Glob(filepath.Join(dir, ".dedao-pdf-*.pdf"))
			if len(staged) != 0 {
				t.Fatal("staged PDF leaked", staged)
			}
		})
	}
}

func TestPDFFileURL(t *testing.T) {
	if got := pdfFileURL("/tmp/中文 space.html"); !strings.HasPrefix(got, "file:///") || strings.Contains(got, " ") {
		t.Fatal(got)
	}
}
