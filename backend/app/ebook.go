package app

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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
	wgp := utils.NewWaitGroupPool(5)
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
			svgList, err1 := generateEbookPages(enID, order.ChapterID, token.Token, index, count, offset)
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

func generateEbookPages(enid, chapterID, token string, index, count, offset int) (svgList []string, err error) {
	// Try to load from cache first
	if cachedPages, found := services.LoadFromCache(enid, chapterID); found {
		fmt.Printf("使用缓存内容：%s\n", chapterID)
		return cachedPages, nil
	}

	fmt.Printf("下载章节 %s\n", chapterID)
	pageList, err := getService().EbookPages(chapterID, token, index, count, offset)
	if err != nil {
		return
	}

	for _, item := range pageList.Pages {
		desContents := DecryptAES(item.Svg)
		svgList = append(svgList, desContents)
	}

	if !pageList.IsEnd {
		index += count
		count = 20
		fmt.Printf("下载章节 %s 的更多页面 (索引: %d)\n", chapterID, index)
		list, err1 := generateEbookPages(enid, chapterID, token, index, count, offset)
		if err1 != nil {
			err = err1
			return
		}

		svgList = append(svgList, list...)
	} else {
		fmt.Printf("章节 %s 下载完成 (共 %d 页)\n", chapterID, len(svgList))
	}

	// Save to cache
	if err := services.SaveToCache(enid, chapterID, svgList); err != nil {
		fmt.Printf("警告: 无法缓存章节 %s: %v\n", chapterID, err)
	} else {
		fmt.Printf("已缓存章节 %s\n", chapterID)
	}

	return
}

// PKCS7Unpad 实现PKCS7去填充
func PKCS7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// DecryptAES 实现AES - CBC解密
func DecryptAES(contents string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(contents)
	if err != nil {
		fmt.Println("Base64解码错误:", err)
		return ""
	}

	key := []byte("3e4r06tjkpjcevlbslr3d96gdb5ahbmo")
	iv := []byte("6fd89a1b3a7f48fb")

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	blockSize := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = PKCS7Unpad(plaintext)
	return string(plaintext)
}
