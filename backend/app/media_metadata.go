package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yann0917/dedao-gui/backend/downloader"
	"github.com/yann0917/dedao-gui/backend/mediameta"
	"github.com/yann0917/dedao-gui/backend/services"
	"github.com/yann0917/dedao-gui/backend/utils"
)

const DownloadMetadata = 4

func exportAudioMetadata(ctx context.Context, dir string, album mediameta.Album, tracks []mediameta.Track) error {
	if ctx == nil {
		ctx = context.Background()
	}
	// Composition root: portable metadata export depends only on HTTP/tagger
	// interfaces; the desktop setting chooses the concrete ffmpeg executable.
	exporter := mediameta.Exporter{Tagger: mediameta.FFmpegTagger{Executable: utils.FfmpegDir}}
	return exporter.Export(ctx, dir, album, tracks)
}

func courseMetadata(course *services.CourseInfo, articles *services.ArticleList, aid int) (mediameta.Album, []mediameta.Track) {
	c := course.ClassInfo
	album := mediameta.Album{
		Title: c.Name, Artist: c.LecturerName, ID: strconv.Itoa(c.ID),
		Description: plainMetadataText(firstMetadataText(c.Intro, c.ShareSummary, c.Highlight)),
		Date:        metadataDate(c.PublishTime),
		CoverURL:    firstMetadataText(c.SquareImg, c.Logo, c.IndexImg, c.PlayerImg),
	}
	tracks := make([]mediameta.Track, 0)
	for i, article := range articles.List {
		if (aid > 0 && article.ID != aid) || article.Audio == nil {
			continue
		}
		number := article.OrderNum
		if number <= 0 {
			number = i + 1
		}
		tracks = append(tracks, mediameta.Track{
			Filename: downloader.AudioFilename(downloader.Datum{Title: article.Title, OrderNum: article.OrderNum}),
			Title:    article.Title, Number: number,
		})
	}
	return album, tracks
}

func odobMetadata(course *services.Course) (mediameta.Album, []mediameta.Track) {
	return mediameta.Album{
		Title: course.Title, Artist: course.Author, ID: strconv.Itoa(course.ID),
		Description: plainMetadataText(firstMetadataText(course.Intro, course.ProductIntro)),
		CoverURL:    course.Icon,
	}, []mediameta.Track{{Filename: downloader.AudioFilename(downloader.Datum{Title: course.Title}), Title: course.Title, Number: 1}}
}

// Every book needs its own album/cover; IDs also disambiguate equal titles.
// Move a matching legacy flat MP3 only when the new destination is absent.
func odobAudioDir(output string, course *services.Course) (string, error) {
	root := filepath.Join(output, "每天听本书", "MP3")
	dir, err := utils.Mkdir(root, fmt.Sprintf("%s [%d]", utils.FileName(course.Title, ""), course.ID))
	if err != nil {
		return "", err
	}
	name := downloader.AudioFilename(downloader.Datum{Title: course.Title})
	target := filepath.Join(dir, name)
	if _, err = os.Stat(target); os.IsNotExist(err) {
		legacy := filepath.Join(root, name)
		if info, statErr := os.Stat(legacy); statErr == nil && info.Mode().IsRegular() {
			if err = os.Rename(legacy, target); err != nil {
				return "", err
			}
		} else if statErr != nil && !os.IsNotExist(statErr) {
			return "", statErr
		}
	} else if err != nil {
		return "", err
	}
	return dir, nil
}

func metadataDate(timestamp int) string {
	if timestamp <= 0 {
		return ""
	}
	return time.Unix(int64(timestamp), 0).UTC().Format("2006-01-02")
}

func firstMetadataText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func plainMetadataText(value string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(value))
	if err != nil {
		return value
	}
	doc.Find("script,style").Remove()
	return strings.TrimSpace(doc.Text())
}
