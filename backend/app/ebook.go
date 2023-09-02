package app

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao-gui/backend/services"
	"github.com/yann0917/dedao-gui/backend/utils"
)

func EbookDetail(enID string) (detail *services.EbookDetail, err error) {
	detail, err = getService().EbookDetail(enID)
	return
}

// EbookCommentList get ebook 评分&书评
// sort like_count
func EbookCommentList(enID, sort string, page, limit int) (list *services.EbookCommentList, err error) {
	list, err = getService().EbookCommentList(enID, sort, page, limit)
	return
}

// EbookShelfAdd 加入书架
func EbookShelfAdd(enIDs []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = getService().EbookShelfAdd(enIDs)
	return
}

// EbookShelfRemove 移出书架
func EbookShelfRemove(enIDs []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = getService().EbookShelfRemove(enIDs)
	return
}

func EbookInfo(enID string) (info *services.EbookInfo, err error) {
	token, err1 := getService().EbookReadToken(enID)
	if err1 != nil {
		err = err1
		return
	}

	info, err = getService().EbookInfo(token.Token)
	return
}

func EbookPage(ctx context.Context, enID string) (info *services.EbookInfo, svgContent utils.SvgContents, err error) {
	token, err1 := getService().EbookReadToken(enID)
	if err1 != nil {
		err = err1
		return
	}

	info, err = getService().EbookInfo(token.Token)
	if err != nil {
		return
	}
	wgp := utils.NewWaitGroupPool(10)
	total, curr := len(info.BookInfo.Orders), 0
	var chapterMap sync.Map
	for _, ebookToc := range info.BookInfo.Toc {
		key := ebookToc.Href
		href := strings.Split(ebookToc.Href, "#")
		if len(href) > 1 {
			key = href[0]
		}
		chapterMap.Store(key, ebookToc)
	}
	for i, order := range info.BookInfo.Orders {
		var progress Progress
		progress.Total = total
		curr++
		progress.Current = curr
		progress.Pct = curr * 100 / progress.Total
		value, ok := chapterMap.Load(order.ChapterID)
		if ok {
			progress.Value = value.(utils.EbookToc).Text
			chapterMap.Delete(order.ChapterID)
		}
		runtime.EventsEmit(ctx, "ebookDownload", progress)
		wgp.Add()
		go func(i int, order services.EbookOrders) {
			defer func() {
				wgp.Done()
			}()
			index, count, offset := 0, 20, 0
			svgList, err1 := generateEbookPages(order.ChapterID, token.Token, index, count, offset)
			if err1 != nil {
				err = err1
				return
			}

			svgContent = append(svgContent, &utils.SvgContent{
				Contents:   svgList,
				ChapterID:  order.ChapterID,
				PathInEpub: order.PathInEpub,
				OrderIndex: i,
			})
		}(i, order)
	}
	wgp.Wait()
	return
}

func generateEbookPages(chapterID, token string, index, count, offset int) (svgList []string, err error) {
	fmt.Printf("chapterID:%#v\n", chapterID)
	pageList, err := getService().EbookPages(chapterID, token, index, count, offset)
	if err != nil {
		return
	}

	for _, item := range pageList.Pages {
		svgList = append(svgList, item.Svg)
	}
	// fmt.Printf("IsEnd:%#v\n", pageList.IsEnd)
	if !pageList.IsEnd {
		index = count
		count += 20
		list, err1 := generateEbookPages(chapterID, token, index, count, offset)
		if err1 != nil {
			err = err1
			return
		}

		svgList = append(svgList, list...)
	}
	// FIXME: debug
	// err = utils.SaveFile(OutputDir, chapterID, "", strings.Join(svgList, "\n"))
	return
}
