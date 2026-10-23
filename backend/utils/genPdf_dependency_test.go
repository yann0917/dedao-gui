package utils

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func TestGenPdfReportsMissingConverter(t *testing.T) {
	t.Setenv("DEDAO_PDF_ENGINE", "wkhtmltopdf")
	old := WkToPdfDir
	WkToPdfDir = filepath.Join(t.TempDir(), "missing-wkhtmltopdf")
	t.Cleanup(func() {
		WkToPdfDir = old
		wkhtmltopdf.SetPath("")
	})
	p := PdfOption{FileName: filepath.Join(t.TempDir(), "out.pdf")}
	err := p.GenPdf(bytes.NewBufferString("<html><body>test</body></html>"))
	if err == nil || !strings.Contains(err.Error(), "wkhtmltopdf") || !strings.Contains(err.Error(), "系统设置") {
		t.Fatalf("expected actionable missing wkhtmltopdf error, got %v", err)
	}
}
