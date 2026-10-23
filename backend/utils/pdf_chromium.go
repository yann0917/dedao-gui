package utils

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/net/html"
)

var ChromiumPath string

type chromiumPDFRenderer struct{}

func findChromium() (string, error) {
	configured := ChromiumPath
	if configured == "" {
		configured = os.Getenv("DEDAO_CHROME_PATH")
	}
	return locateChromium(configured, chromiumCandidates())
}

func locateChromium(configured string, candidates []string) (string, error) {
	if configured != "" {
		return findExecutable(configured, "Chrome/Edge/Chromium", nil)
	}
	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到 Chrome、Edge 或 Chromium，请先安装浏览器，或在系统设置中填写 PDF 浏览器路径（也可设置 DEDAO_CHROME_PATH）")
}

func chromiumCandidates() []string {
	names := []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "microsoft-edge", "msedge", "chrome"}
	switch runtime.GOOS {
	case "darwin":
		for _, root := range []string{"/Applications", filepath.Join(os.Getenv("HOME"), "Applications")} {
			for _, app := range []string{"Google Chrome", "Microsoft Edge", "Chromium", "Google Chrome for Testing"} {
				names = append(names, filepath.Join(root, app+".app", "Contents", "MacOS", app))
			}
		}
	case "windows":
		for _, root := range []string{os.Getenv("PROGRAMFILES"), os.Getenv("PROGRAMFILES(X86)"), os.Getenv("LOCALAPPDATA")} {
			if root == "" {
				continue
			}
			names = append(names, filepath.Join(root, "Google", "Chrome", "Application", "chrome.exe"), filepath.Join(root, "Microsoft", "Edge", "Application", "msedge.exe"), filepath.Join(root, "Chromium", "Application", "chrome.exe"))
		}
	}
	return names
}

func (chromiumPDFRenderer) Render(ctx context.Context, p *PdfOption, data []byte) error {
	browser, err := findChromium()
	if err != nil {
		return err
	}
	if p.PageSize != "" && p.PageSize != "A4" {
		return fmt.Errorf("Chromium PDF 暂仅支持 A4 纸张")
	}
	var cover []byte
	if p.CoverPath != "" {
		cover, err = os.ReadFile(p.CoverPath)
		if err != nil {
			return fmt.Errorf("读取 PDF 封面失败: %w", err)
		}
	}
	document, err := chromiumDocument(data, cover)
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "dedao-pdf-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	input := filepath.Join(dir, "document.html")
	if err = os.WriteFile(input, document, 0600); err != nil {
		return err
	}
	// Stage beside the target: a failed browser must not truncate an existing PDF.
	output, err := os.CreateTemp(filepath.Dir(p.FileName), ".dedao-pdf-*.pdf")
	if err != nil {
		return err
	}
	if err = output.Close(); err != nil {
		os.Remove(output.Name())
		return err
	}
	defer os.Remove(output.Name())
	outPath, err := filepath.Abs(output.Name())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	inputURL := pdfFileURL(input)
	args := []string{
		"--headless", "--no-first-run", "--no-default-browser-check", "--disable-extensions",
		"--disable-background-networking", "--disable-sync", "--disable-component-update",
		"--user-data-dir=" + filepath.Join(dir, "profile"),
		"--no-pdf-header-footer", "--generate-pdf-document-outline",
		"--timeout=60000", "--virtual-time-budget=1000", "--print-to-pdf=" + outPath, inputURL,
	}
	cmd := exec.CommandContext(ctx, browser, args...)
	configurePDFProcess(cmd)
	cmd.WaitDelay = 5 * time.Second
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("PDF 浏览器渲染超时或已取消: %w", ctx.Err())
		}
		message := stderr.String()
		if len(message) > 2000 {
			message = message[len(message)-2000:]
		}
		return fmt.Errorf("Chromium PDF 渲染失败: %w\n%s", err, message)
	}
	f, err := os.Open(outPath)
	if err != nil {
		return fmt.Errorf("浏览器未生成 PDF: %w", err)
	}
	header := make([]byte, 5)
	_, readErr := f.Read(header)
	f.Close()
	if readErr != nil || string(header) != "%PDF-" {
		return fmt.Errorf("浏览器未生成有效 PDF，请检查浏览器路径和版本")
	}
	if pages, err := api.PageCountFile(ctx, outPath); err != nil || pages < 1 {
		return fmt.Errorf("浏览器生成的 PDF 无有效页面: %v", err)
	}
	return os.Rename(outPath, p.FileName)
}

func pdfFileURL(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return (&url.URL{Scheme: "file", Path: path}).String()
}

// Merge fragments into one print document so cover pages cannot shift links or
// duplicate outlines during a later PDF merge. Never run scripts from content.
func chromiumDocument(body, cover []byte) ([]byte, error) {
	styles, content, err := printableFragment(body)
	if err != nil {
		return nil, err
	}
	var result strings.Builder
	result.WriteString(`<!DOCTYPE html><html lang="zh-CN"><head><meta charset="UTF-8"><meta http-equiv="Content-Security-Policy" content="default-src 'none'; script-src 'none'; style-src 'unsafe-inline'; img-src file: data: https: http:; font-src file: data: https: http:; base-uri 'none'; form-action 'none'">`)
	result.WriteString(styles)
	result.WriteString(`<style>
@page { size:A4; margin:15mm; @bottom-right { content:counter(page); font:10pt sans-serif; } }
@page dedao-cover { size:A4; margin:0; @bottom-right { content:none; } }
html,body { margin:0; padding:0; }
body { font-family:"PingFang SC","Microsoft YaHei","Noto Sans CJK SC",sans-serif; line-height:1.8; color:#333; }
* { -webkit-print-color-adjust:exact; print-color-adjust:exact; }
h1,h2,h3,h4,h5,h6 { break-after:avoid; }
img { max-width:100%; }
.dedao-pdf-cover { page:dedao-cover; width:210mm; height:297mm; overflow:hidden; break-after:page; }
.dedao-pdf-cover img { display:block; width:210mm; height:297mm; max-width:none!important; object-fit:fill; }
.dedao-pdf-main { page:auto; }
</style></head><body>`)
	if len(cover) != 0 {
		_, coverContent, err := printableFragment(cover)
		if err != nil {
			return nil, err
		}
		result.WriteString(`<section class="dedao-pdf-cover">` + coverContent + `</section>`)
	}
	result.WriteString(`<main class="dedao-pdf-main">` + content + `</main></body></html>`)
	return []byte(result.String()), nil
}

func printableFragment(data []byte) (string, string, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return "", "", fmt.Errorf("解析 PDF HTML 失败: %w", err)
	}
	removeActivePDFContent(doc)
	var styles, body bytes.Buffer
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "head" || n.Data == "body") {
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if n.Data == "body" {
					html.Render(&body, child)
				} else if child.Type == html.ElementNode && child.Data == "style" {
					html.Render(&styles, child)
				}
			}
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)
	return styles.String(), body.String(), nil
}

func removeActivePDFContent(n *html.Node) {
	for child := n.FirstChild; child != nil; {
		next := child.NextSibling
		blocked := false
		if child.Type == html.ElementNode {
			switch child.Data {
			case "script", "iframe", "object", "embed", "meta", "link", "base":
				blocked = true
			}
			attrs := child.Attr[:0]
			for _, attr := range child.Attr {
				if !strings.HasPrefix(strings.ToLower(attr.Key), "on") {
					attrs = append(attrs, attr)
				}
			}
			child.Attr = attrs
		}
		if blocked {
			n.RemoveChild(child)
		} else {
			removeActivePDFContent(child)
		}
		child = next
	}
}

func hasVisiblePDFContent(data []byte) bool {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return false
	}
	var visible func(*html.Node, bool) bool
	visible = func(n *html.Node, inBody bool) bool {
		if n.Type == html.ElementNode {
			if n.Data == "head" || n.Data == "style" || n.Data == "script" {
				return false
			}
			inBody = inBody || n.Data == "body"
			if inBody && (n.Data == "img" || n.Data == "svg" || n.Data == "math") {
				return true
			}
		}
		if inBody && n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			return true
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if visible(child, inBody) {
				return true
			}
		}
		return false
	}
	return visible(doc, false)
}
