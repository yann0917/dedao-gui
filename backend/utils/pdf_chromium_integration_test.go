package utils

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestGenPdfWithoutWkhtmltopdf(t *testing.T) {
	if os.Getenv("DEDAO_TEST_CHROMIUM") == "" {
		t.Skip("set DEDAO_TEST_CHROMIUM=1 for real browser integration")
	}
	t.Setenv("DEDAO_PDF_ENGINE", "chromium")
	old := WkToPdfDir
	WkToPdfDir = filepath.Join(t.TempDir(), "no-wkhtmltopdf")
	t.Cleanup(func() { WkToPdfDir = old })
	dir := t.TempDir()
	if out := os.Getenv("DEDAO_PDF_TEST_OUTPUT"); out != "" {
		dir = out
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	cover := filepath.Join(dir, "cover.html")
	if err := os.WriteFile(cover, []byte(`<html><body><h1>中文封面</h1><p>现代 PDF 导出验证</p></body></html>`), 0600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "modern-pdf-smoke.pdf")
	p := PdfOption{FileName: out, CoverPath: cover, PageSize: "A4", Toc: true}
	markup := `<html><head><meta charset="UTF-8"></head><body><h1>目 录</h1><a href="#chapter">第一章：中文排版</a><section style="break-before:page"><h1 id="chapter">第一章：中文排版</h1><p>这是中文正文，验证分页、目录链接和字体。</p><p>English paragraph with punctuation &amp; symbols.</p></section></body></html>`
	if err := p.GenPdf(bytes.NewBufferString(markup)); err != nil {
		t.Fatalf("PDF export must work without wkhtmltopdf: %v", err)
	}
	pages, err := api.PageCountFile(context.Background(), out)
	if err != nil || pages != 3 {
		t.Fatalf("want cover + contents + chapter (3 pages), got %d: %v", pages, err)
	}
	if _, err := exec.LookPath("pdftotext"); err == nil {
		text, err := exec.Command("pdftotext", out, "-").Output()
		if err != nil || !strings.Contains(string(text), "这是中文正文") || !strings.Contains(string(text), "中文封面") {
			t.Fatalf("Chinese text not preserved: %s (%v)", text, err)
		}
	}
	cc, err := api.ReadContextFile(context.Background(), out)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cc.RootDict["Outlines"]; !ok {
		t.Fatal("missing document outline")
	}
	assertPDFDestinationPage(t, out, "chapter", 3)
	t.Logf("Verified PDF: %s (%d pages)", out, pages)
}

func TestSvg2PdfChromium(t *testing.T) {
	if os.Getenv("DEDAO_TEST_CHROMIUM") == "" {
		t.Skip("set DEDAO_TEST_CHROMIUM=1 for real browser integration")
	}
	t.Setenv("DEDAO_PDF_ENGINE", "chromium")
	dir := t.TempDir()
	if out := os.Getenv("DEDAO_PDF_TEST_OUTPUT"); out != "" {
		dir = out
	}
	coverImage := image.NewRGBA(image.Rect(0, 0, 2, 3))
	for y := 0; y < 3; y++ {
		for x := 0; x < 2; x++ {
			coverImage.Set(x, y, color.RGBA{R: 160, G: 25, B: 35, A: 255})
		}
	}
	coverPath := filepath.Join(dir, "image-cover.png")
	f, err := os.Create(coverPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, coverImage); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	pages := SvgContents{
		{ChapterID: "cover.xhtml", OrderIndex: 0, Contents: []string{fixtureSvgPage(`<image x="0" y="0" width="60000" height="85000" href="` + pdfFileURL(coverPath) + `"/>`)}},
		{ChapterID: "chapter1.xhtml", OrderIndex: 1, Contents: []string{fixtureSvgPage(fixtureTextLine("first", "0", "100", "第一章 中文测试", "22px"), fixtureTextLine("", "10", "200", "正文段落内容。", "16px"))}},
		{ChapterID: "chapter2.xhtml", OrderIndex: 2, Contents: []string{fixtureSvgPage(fixtureTextLine("second", "20", "100", "第二章 分页测试", "22px"), fixtureTextLine("", "30", "200", "第二章正文。", "16px"))}},
	}
	toc := []EbookToc{{Href: "chapter1.xhtml#first", Text: "第一章 中文测试", Level: 1}, {Href: "chapter2.xhtml#second", Text: "第二章 分页测试", Level: 1}}
	if err := Svg2Pdf(dir, "电子书渲染验证", pages, toc); err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "Ebook", "*.pdf"))
	if len(files) == 0 {
		t.Fatal("no ebook PDF")
	}
	out := files[len(files)-1]
	count, err := api.PageCountFile(context.Background(), out)
	if err != nil || count != 4 {
		t.Fatalf("want cover + contents + 2 chapters, got %d: %v (%s)", count, err, out)
	}
	assertPDFDestinationPage(t, out, "first", 3)
	assertPDFDestinationPage(t, out, "second", 4)
}

func assertPDFDestinationPage(t *testing.T, path, name string, want int) {
	t.Helper()
	cc, err := api.ReadContextFile(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	dests, err := cc.DereferenceDict(cc.RootDict["Dests"])
	if err != nil {
		t.Fatal(err)
	}
	dest, err := cc.DereferenceArray(dests[name])
	if err != nil || len(dest) == 0 {
		t.Fatalf("missing destination %s: %v", name, err)
	}
	ref, ok := dest[0].(types.IndirectRef)
	if !ok {
		t.Fatalf("destination %s is not a page reference: %v", name, dest)
	}
	got, err := cc.PageNumber(context.Background(), ref.ObjectNumber.Value())
	if err != nil || got != want {
		t.Fatalf("destination %s points to %d, want %d: %v", name, got, want, err)
	}
}
