package app

import (
	"fmt"
	"math"

	"github.com/yann0917/dedao-gui/backend/services"
)

// ArticleList 已购课程文章列表分页
func ArticleList(enid, chapterID string, count, maxID int) (list *services.ArticleList, err error) {
	list, err = getService().ArticleList(enid, chapterID, count, maxID)
	if err != nil {
		return
	}
	return
}

// ArticleListAll 已购课程文章列表
func ArticleListAll(id int, chapterID string) (list *services.ArticleList, err error) {
	info, err := CourseDetail(CateCourse, id)
	if err != nil {
		return
	}
	pageSize := 30
	enid := info.Enid
	count := info.PublishNum
	page := int(math.Ceil(float64(count) / float64(pageSize)))
	maxID := 0
	var lists []services.ArticleIntro
	for i := 0; i < page; i++ {
		list, err = getService().ArticleList(enid, chapterID, pageSize, maxID)
		if err != nil {
			return
		}
		maxID = list.List[len(list.List)-1].ID
		lists = append(lists, list.List...)
	}
	list.List = lists
	return

}

// ArticleDetail article detail
func ArticleDetail(enid string) (detail *services.ArticleDetail, err error) {
	info, err := getService().ArticleInfo(enid, 1)
	if err != nil {
		return
	}
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, enid, appid)
	if err != nil {
		return
	}
	return

}

// ArticleCommentList article comment list
func ArticleCommentList(enId, sort string, page, limit, sType int) (list *services.ArticleCommentList, err error) {
	list, err = getService().ArticleCommentList(enId, sort, page, limit, sType)
	if err != nil {
		return
	}

	return

}

// OdobArticleInfo article info
// get article token, audio token, media security token etc.
func OdobArticleInfo(aEnid string) (info *services.ArticleInfo, err error) {
	info, err = getService().ArticleInfo(aEnid, 2)
	if err != nil {
		return
	}
	return
}

// OdobArticleDetail odob article detail
func OdobArticleDetail(aEnid string) (detail *services.ArticleDetail, err error) {
	info, err := OdobArticleInfo(aEnid)
	if err != nil {
		return
	}

	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, aEnid, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}

// ArticleDetailByEnid odob article detail
// enid article enid  or odob audioAliasID, aType 1-course article, 2-odob article
func ArticleDetailByEnid(aType int, aEnid string) (detail *services.ArticleDetail, err error) {
	info, err := getService().ArticleInfo(aEnid, aType)
	if err != nil {
		return
	}

	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, aEnid, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}
