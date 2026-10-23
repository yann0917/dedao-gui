package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yann0917/dedao-gui/backend/downloader"
	"github.com/yann0917/dedao-gui/backend/services"
)

func TestCourseMetadataMappingAndArticleSelection(t *testing.T) {
	course := &services.CourseInfo{ClassInfo: services.ClassInfo{
		ID: 123, Name: "中文课程", LecturerName: "讲师", Intro: "<p>简介 &amp; 文明</p><script>ignored()</script>",
		Logo: "https://example.com/logo", SquareImg: "https://example.com/square", PublishTime: 1792713600,
	}}
	articles := &services.ArticleList{List: []services.ArticleIntro{
		{ArticleBase: services.ArticleBase{ID: 1, Title: "长标题/" + strings.Repeat("文", 90), OrderNum: 12}, Audio: &services.Audio{}},
		{ArticleBase: services.ArticleBase{ID: 2, Title: "没有音频"}},
		{ArticleBase: services.ArticleBase{ID: 3, Title: "未编号"}, Audio: &services.Audio{}},
	}}
	album, tracks := courseMetadata(course, articles, 1)
	if album.Title != "中文课程" || album.Artist != "讲师" || album.ID != "123" || album.Description != "简介 & 文明" || album.CoverURL != course.ClassInfo.SquareImg || album.Date != "2026-10-23" {
		t.Fatalf("unexpected album: %+v", album)
	}
	if len(tracks) != 1 || tracks[0].Number != 12 || tracks[0].Filename != downloader.AudioFilename(downloader.Datum{Title: articles.List[0].Title, OrderNum: 12}) {
		t.Fatalf("single-article selection/naming mismatch: %+v", tracks)
	}
	_, tracks = courseMetadata(course, articles, 0)
	if len(tracks) != 2 || tracks[1].Number != 3 || tracks[1].Filename != "未编号.mp3" {
		t.Fatalf("unexpected fallback track number: %+v", tracks)
	}
}

func TestOdobDirectoriesAreIsolatedAndMigrateLegacyMP3(t *testing.T) {
	output := t.TempDir()
	book := &services.Course{ID: 42, Title: "书/同名", Author: "作者", Intro: "<p>书籍简介</p>"}
	name := downloader.AudioFilename(downloader.Datum{Title: book.Title})
	root := filepath.Join(output, "每天听本书", "MP3")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(root, name)
	if err := os.WriteFile(legacy, []byte("legacy audio"), 0644); err != nil {
		t.Fatal(err)
	}
	dir, err := odobAudioDir(output, book)
	if err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(dir, name)); err != nil || string(data) != "legacy audio" {
		t.Fatalf("legacy audio not moved: %s %v", data, err)
	}
	// If both layouts exist, never overwrite/delete either user's audio file.
	if err := os.WriteFile(legacy, []byte("another copy"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err = odobAudioDir(output, book); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(legacy); string(data) != "another copy" {
		t.Fatal("existing target should preserve legacy copy")
	}
	book.ID = 43
	other, err := odobAudioDir(output, book)
	if err != nil || dir == other {
		t.Fatalf("same titles must not share covers/NFO: %s %v", other, err)
	}
	album, tracks := odobMetadata(book)
	if album.Date != "" || album.Description != "书籍简介" || album.Artist != "作者" || len(tracks) != 1 || tracks[0].Filename != name {
		t.Fatalf("unexpected audiobook metadata: %+v %+v", album, tracks)
	}
}

func TestOdobMetadataOnlyDoesNotFetchStreamsOrRequireLogin(t *testing.T) {
	old := OutputDir
	OutputDir = t.TempDir()
	t.Cleanup(func() { OutputDir = old })
	// No account, alias ID, stream URL, or desktop application is configured.
	d := OdobDownload{Ctx: context.Background(), DownloadType: DownloadMetadata, Data: &services.Course{ID: 10, Title: "只导出元数据", Author: "讲书人"}}
	if err := d.Download(); err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join(OutputDir, "每天听本书", "MP3", "*", "album.nfo"))
	if len(files) != 1 {
		t.Fatalf("metadata-only task did not write NFO: %v", files)
	}
	if err := (&OdobDownload{DownloadType: DownloadMetadata}).Download(); err == nil {
		t.Fatal("nil audiobook data should return an error, not panic")
	}
}
