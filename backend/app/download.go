package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao-gui/backend/downloader"
	"github.com/yann0917/dedao-gui/backend/services"
	"github.com/yann0917/dedao-gui/backend/utils"
)

var OutputDir = ""

type CourseDownload struct {
	Ctx          context.Context
	DownloadType int    // 1:mp3, 2:PDF文档, 3:markdown文档
	ID           int    // 课程 id
	AID          int    // 文章 id
	EnId         string // 课程 enid
}

type OdobDownload struct {
	Ctx          context.Context
	DownloadType int // 1:mp3, 2:PDF文档, 3:markdown文档
	ID           int
	Data         *services.Course
}

type EBookDownload struct {
	Ctx          context.Context
	DownloadType int // 1:html, 2:PDF文档, 3:epub
	ID           int
	EnID         string
}

type Progress struct {
	Total   int    `json:"total"`
	Current int    `json:"current"`
	Pct     int    `json:"pct"`
	Value   string `json:"value"`
	ID      int    `json:"id"` // 课程 id
}

func SetOutputDir(dir string) {
	OutputDir = dir
}

func (d *CourseDownload) Download() error {
	course, err := getService().CourseInfo(d.EnId)
	if err != nil {
		return err
	}
	count := course.ClassInfo.CurrentArticleCount
	if count == 0 {
		count = course.ClassInfo.FormalArticleCount
	}
	// Fallback to ensure we at least try to fetch one page
	if count == 0 {
		count = 100
	}

	articles, err := ArticleListAllByInfo(d.EnId, count, "")
	if err != nil {
		return err
	}

	switch d.DownloadType {
	case 1: // mp3
		downloadData := extractDownloadData(course, articles, d.AID, 1)
		errs := make([]error, 0)

		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "MP3")
		if err != nil {
			return err
		}

		total, curr := len(downloadData.Data), 0
		for _, datum := range downloadData.Data {
			var progress Progress
			progress.ID = d.ID
			progress.Total = total
			curr++
			progress.Current = curr
			progress.Pct = curr * 100 / progress.Total
			progress.Value = datum.Title
			runtime.EventsEmit(d.Ctx, "courseDownload", progress)
			if !datum.IsCanDL {
				continue
			}
			stream := datum.Enid
			if err := downloader.Download(datum, stream, path); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errs[0]
		}
	case 2:
		// 下载 PDF
		downloadData := extractDownloadData(course, articles, d.AID, 2)

		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "PDF")
		if err != nil {
			return err
		}
		return DownloadPdfCourse(downloadData.Data, path, d.Ctx)

	case 3:
		// 下载 Markdown
		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "MD")
		if err != nil {
			return err
		}
		return DownloadMarkdown(articles, d.AID, path, d.Ctx)
	}
	return nil

}

func (d *OdobDownload) Download() error {
	fileName := "每天听本书"
	switch d.DownloadType {
	case 1:
		downloadData := downloader.Data{
			Title: fileName,
		}
		downloadData.Type = "audio"
		downloadData.Data = extOdobDownloadData(d.Data)
		errors := make([]error, 0)
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "MP3")
		if err != nil {
			return err
		}
		total, curr := len(downloadData.Data), 0
		for _, datum := range downloadData.Data {
			var progress Progress
			progress.ID = d.ID
			progress.Total = total
			curr++
			progress.Current = curr
			progress.Pct = curr * 100 / progress.Total
			progress.Value = datum.Title + ".mp3"
			runtime.EventsEmit(d.Ctx, "odobDownload", progress)
			if !datum.IsCanDL {
				continue
			}
			stream := datum.Enid
			if err := downloader.Download(datum, stream, path); err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return errors[0]
		}
	case 2:
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "PDF")
		if err != nil {
			return err
		}
		aliasID := d.Data.AudioDetail.AliasID
		detail, err := OdobArticleDetail(aliasID)
		if err != nil {
			return err
		}

		var content []services.Content
		err = jsoniter.UnmarshalFromString(detail.Content, &content)
		if err != nil {
			return err
		}
		res := ContentsToMarkdown(content)
		var progress Progress
		progress.ID = d.ID
		progress.Total = 100
		progress.Current = 100
		progress.Pct = 100 * 100 / progress.Total
		progress.Value = d.Data.Title + ".pdf"
		runtime.EventsEmit(d.Ctx, "odobDownload", progress)
		return utils.Md2Pdf(path, d.Data.Title, []byte(res))
	case 3:
		// 下载 Markdown
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "MD")
		if err != nil {
			return err
		}
		var progress Progress
		progress.ID = d.ID
		progress.Total = 100
		progress.Current = 100
		progress.Pct = 100 * 100 / progress.Total
		progress.Value = d.Data.Title + ".md"
		runtime.EventsEmit(d.Ctx, "odobDownload", progress)
		if err := DownloadOdobMarkdown(d.Data, path); err != nil {
			return err
		}
	}
	return nil
}

func (d *EBookDownload) Download() error {
	detail, err := EbookDetail(d.EnID)
	if err != nil {
		return err
	}

	title := strconv.Itoa(d.ID) + "_"
	if detail.Title != "" {
		title += detail.Title
	} else if detail.OperatingTitle != "" {
		title += detail.OperatingTitle
	}

	title += "_" + detail.BookAuthor
	info, svgContent, err := EbookPage(d.Ctx, detail.Enid)
	if err != nil {
		return err
	}
	sort.Sort(svgContent)

	dType := map[int]string{
		1: "HTML",
		2: "PDF",
		3: "EPUB",
	}
	var progress Progress
	progress.Pct = 100
	progress.Value = "正在生成" + dType[d.DownloadType] + "文件"
	runtime.EventsEmit(d.Ctx, "ebookDownload", progress)
	switch d.DownloadType {
	case 1:
		if err = utils.Svg2Html(OutputDir, title, svgContent, info.BookInfo.Toc); err != nil {
			return err
		}

	case 2:
		if err = utils.Svg2Pdf(OutputDir, title, svgContent, info.BookInfo.Toc); err != nil {
			return err
		}

	case 3:
		var opts utils.EpubOptions
		opts.Title = title
		opts.Author = detail.BookAuthor
		opts.Description = detail.BookIntro
		opts.Toc = info.BookInfo.Toc

		if err = utils.Svg2Epub(OutputDir, title, svgContent, opts); err != nil {
			return err
		}

		return err
	}

	return nil
}

// 生成下载数据
func extractDownloadData(course *services.CourseInfo, articles *services.ArticleList, aid int, flag int) downloader.Data {

	downloadData := downloader.Data{
		Title: course.ClassInfo.Name,
	}

	if course.HasAudio() {
		downloadData.Type = "audio"
		downloadData.Data = extractCourseDownloadData(articles, aid, flag)
	}

	return downloadData
}

// 生成课程下载数据
func extractCourseDownloadData(articles *services.ArticleList, aid int, flag int) []downloader.Datum {
	data := downloader.EmptyData
	audioIds := map[int]string{}

	audioData := make([]*downloader.Datum, 0)
	for _, article := range articles.List {
		if aid > 0 && article.ID != aid {
			continue
		}

		if article.Audio.Mp3PlayURL != "" && len(article.AudioAliasIds) > 0 {
			audioIds[article.ID] = article.Audio.AliasID

			var urls []downloader.URL
			key := article.Enid
			streams := map[string]downloader.Stream{
				key: {
					URLs:    urls,
					Size:    article.Audio.Size,
					Quality: key,
				},
			}
			isCanDL := true
			if len(article.Audio.AliasID) == 0 {
				isCanDL = false
			}
			datum := &downloader.Datum{
				ID:        article.ID,
				Enid:      article.Enid,
				ClassEnid: article.ClassEnid,
				ClassID:   article.ClassID,
				Title:     article.Title,
				IsCanDL:   isCanDL,
				M3U8URL:   article.Audio.Mp3PlayURL,
				OrderNum:  article.OrderNum,
				Streams:   streams,
				Type:      "audio",
			}

			audioData = append(audioData, datum)
		}

	}

	if flag == 1 {
		handleStreams(audioData, audioIds)
	}

	for _, d := range audioData {
		data = append(data, *d)
	}
	return data
}

// extOdobDownloadData 生成 AudioBook 下载数据
func extOdobDownloadData(info *services.Course) []downloader.Datum {
	data := downloader.EmptyData
	audioIds := map[int]string{}

	audioData := make([]*downloader.Datum, 0)
	aliasID := info.AudioDetail.AliasID

	audioIds[info.ID] = aliasID

	var urls []downloader.URL
	key := info.Enid
	streams := map[string]downloader.Stream{
		key: {
			URLs:    urls,
			Size:    info.AudioDetail.Size,
			Quality: key,
		},
	}
	isCanDL := true
	if !info.HasPlayAuth {
		isCanDL = false
	}
	detail, err := getService().AudioDetailAlias(aliasID)
	if err != nil {
		return nil
	}
	datum := &downloader.Datum{
		ID:      info.ID,
		Enid:    info.Enid,
		ClassID: info.ClassID,
		Title:   info.Title,
		IsCanDL: isCanDL,
		M3U8URL: detail.Mp3PlayURL,
		Streams: streams,
		Type:    "audio",
	}

	audioData = append(audioData, datum)
	handleStreams(audioData, audioIds)

	for _, d := range audioData {
		data = append(data, *d)
	}
	return data
}

func handleStreams(audioData []*downloader.Datum, audioIds map[int]string) {
	wgp := utils.NewWaitGroupPool(10)
	for _, datum := range audioData {
		wgp.Add()
		go func(datum *downloader.Datum, streams map[int]string) {
			defer func() {
				wgp.Done()
			}()
			if datum.IsCanDL {
				if urls, err := utils.M3u8URLs(datum.M3U8URL); err == nil {
					key := datum.Enid
					stream := datum.Streams[key]
					for _, url := range urls {
						stream.URLs = append(stream.URLs, downloader.URL{
							URL: url,
							Ext: "ts",
						})
					}
					datum.Streams[key] = stream
				}
				for k, v := range datum.Streams {
					if len(v.URLs) == 0 {
						delete(datum.Streams, k)
					}
				}
			}
		}(datum, audioIds)
	}
	wgp.Wait()
}

func ContentsToMarkdown(contents []services.Content) (res string) {
	for _, content := range contents {
		switch content.Type {
		case "audio":
			title := strings.TrimRight(content.Title, ".mp3")
			res += getMdHeader(1) + title + "\r\n\r\n"
		case "header":
			content.Text = strings.Trim(content.Text, " ")
			if len(content.Text) > 0 {
				res += getMdHeader(content.Level) + content.Text + "\r\n\r\n"
			}
		case "blockquote":
			texts := strings.Split(content.Text, "\n")
			for _, text := range texts {
				res += "> " + text + "\r\n"
				res += "> \r\n"
			}
			res = strings.TrimRight(res, "> \r\n")
			res += "\r\n\r\n"
		case "paragraph":
			resP, err := paragraphToMarkDown(content.Contents)
			if err != nil {
				return
			}
			res += resP
		case "list":
			resL, err := listToMarkdown(content.Contents)
			if err != nil {
				return
			}
			res += resL
		case "elite": // 划重点
			temp := strings.ReplaceAll(content.Text, "\n", "\r\n\r\n")
			res += getMdHeader(2) + "划重点\r\n\r\n" + temp + "\r\n\r\n"

		case "image":
			res += "![" + content.URL + "](" + content.URL + ")" + "\r\n\r\n"
		case "label-group":
			res += getMdHeader(2) + "`" + content.Text + "`" + "\r\n\r\n"
		}
	}

	res += "---\r\n"
	return
}

func paragraphToMarkDown(content interface{}) (res string, err error) {
	tmpJson, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}
	cont := services.Contents{}
	err = jsoniter.Unmarshal(tmpJson, &cont)
	if err != nil {
		return
	}
	for _, item := range cont {
		subContent := strings.Trim(item.Text.Content, " ")
		switch item.Type {
		case "text":
			if item.Text.Bold {
				res += " **" + subContent + "** "
			} else if item.Text.Highlight {
				res += " *" + subContent + "* "
			} else {
				res += subContent
			}
		}
	}
	res = strings.Trim(res, " ")
	res = strings.Trim(res, "\r\n")
	res += "\r\n\r\n"
	return
}

func listToMarkdown(content interface{}) (res string, err error) {
	tmpJson, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}
	var cont []services.Contents
	err = jsoniter.Unmarshal(tmpJson, &cont)
	if err != nil {
		return
	}

	for _, item := range cont {
		for _, item := range item {
			subContent := strings.Trim(item.Text.Content, " ")
			switch item.Type {
			case "text":
				if item.Text.Bold {
					res += "* **" + subContent + "** "
				} else if item.Text.Highlight {
					res += "* *" + subContent + "* "
				} else {
					res += "* " + subContent
				}
			}
		}
		res += "\r\n\r\n"
	}
	return
}

func articleCommentsToMarkdown(contents []services.Comment) (res string) {
	res = getMdHeader(2) + "热门留言\r\n\r\n"
	for _, content := range contents {
		res += content.NotesOwner.Name + "：" + content.Note + "\r\n\r\n"
		if content.CommentReply != "" {
			res += "> " + content.CommentReplyUser.Name + "(" + content.CommentReplyUser.Role + ") 回复：" + content.CommentReply + "\r\n\r\n"
		}
	}
	res += "---\r\n"
	return
}

func getMdHeader(level int) string {
	heads := map[int]string{
		1: "# ",
		2: "## ",
		3: "### ",
		4: "#### ",
		5: "##### ",
		6: "###### ",
	}
	if s, ok := heads[level]; ok {
		return s
	}
	return ""
}

func DownloadPdfCourse(list []downloader.Datum, path string, ctx context.Context) error {
	name, fileName := "", ""

	total, curr := len(list), 0
	for _, v := range list {
		var progress Progress
		progress.ID = v.ID
		progress.Total = total
		curr++
		progress.Current = curr
		progress.Pct = curr * 100 / progress.Total
		progress.Value = v.Title
		runtime.EventsEmit(ctx, "courseDownload", progress)
		detail, err := ArticleDetail(v.Enid)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		var content []services.Content
		err = jsoniter.UnmarshalFromString(detail.Content, &content)
		if err != nil {
			return err
		}

		name = utils.FileName(v.Title, "pdf")
		name = fmt.Sprintf("%03d.%s", v.OrderNum, name)
		fileName = filepath.Join(path, name)
		_, exist, err := utils.FileSize(fileName)
		if err != nil {
			return err
		}

		if exist {
			continue
		}

		res := ContentsToMarkdown(content)

		err = utils.Md2Pdf(path, name, []byte(res))
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadMarkdown(list *services.ArticleList, aid int, path string, ctx context.Context) error {
	total, curr := len(list.List), 0
	for _, v := range list.List {
		var progress Progress
		progress.ID = aid
		progress.Total = total
		curr++
		progress.Current = curr
		progress.Pct = curr * 100 / progress.Total
		progress.Value = v.Title
		runtime.EventsEmit(ctx, "courseDownload", progress)

		if aid > 0 && v.ID != aid {
			continue
		}
		detail, err := ArticleDetail(v.Enid)
		if err != nil {
			return err
		}

		var content []services.Content
		err = jsoniter.UnmarshalFromString(detail.Content, &content)
		if err != nil {
			return err
		}

		name := utils.FileName(v.Title, "md")
		name = fmt.Sprintf("%03d.%s", v.OrderNum, name)
		fileName := filepath.Join(path, name)
		_, exist, err := utils.FileSize(fileName)

		if err != nil {
			return err
		}

		if exist {
			return nil
		}

		res := ContentsToMarkdown(content)
		// 添加留言
		commentList, err := ArticleCommentList(v.Enid, "like", 1, 20, 65)
		if err == nil {
			res += articleCommentsToMarkdown(commentList.List)
		}

		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = f.WriteString(res)
		if err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func DownloadOdobMarkdown(info *services.Course, path string) error {
	aliasID := info.AudioDetail.AliasID
	detail, err := OdobArticleDetail(aliasID)
	if err != nil {
		return err
	}

	var content []services.Content
	err = jsoniter.UnmarshalFromString(detail.Content, &content)
	if err != nil {
		return err
	}

	name := utils.FileName(info.Title, "md")
	fileName := filepath.Join(path, name)
	_, exist, err := utils.FileSize(fileName)

	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	res := ContentsToMarkdown(content)

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(res)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		if err != nil {
			return err
		}
	}

	return nil

}
