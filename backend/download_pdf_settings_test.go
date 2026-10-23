package backend

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/yann0917/dedao-gui/backend/utils"
)

func TestModernPDFSettingsIgnoreStaleLegacyPath(t *testing.T) {
	t.Setenv("DEDAO_PDF_ENGINE", "chromium")
	oldWK, oldChrome, oldFFmpeg := utils.WkToPdfDir, utils.ChromiumPath, utils.FfmpegDir
	t.Cleanup(func() { utils.WkToPdfDir, utils.ChromiumPath, utils.FfmpegDir = oldWK, oldChrome, oldFFmpeg })
	a := new(App)
	missing := filepath.Join(t.TempDir(), "missing-wkhtmltopdf")
	if err := a.SetDir([]string{"", "", missing, ""}); err != nil {
		t.Fatalf("old settings block modern renderer: %v", err)
	}
	if err := a.SetDir([]string{"", "", missing, filepath.Join(t.TempDir(), "missing-chrome")}); err == nil || !strings.Contains(err.Error(), "PDF 浏览器") {
		t.Fatalf("missing browser validation: %v", err)
	}
	t.Setenv("DEDAO_PDF_ENGINE", "wkhtmltopdf")
	if err := a.SetDir([]string{"", "", missing}); err == nil || !strings.Contains(err.Error(), "wkhtmltopdf") {
		t.Fatalf("legacy validation bypassed: %v", err)
	}
}
