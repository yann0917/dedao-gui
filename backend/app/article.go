package app

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/yann0917/dedao-gui/backend/services"
	"github.com/yann0917/dedao-gui/backend/utils"
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
	enid := info["enid"].(string)
	count := info["publish_num"].(int)
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

// ArticleInfo article info
// get article token, audio token, media security token etc
func ArticleInfo(id, aid int) (info *services.ArticleInfo, aEnid string, err error) {
	list, err := ArticleListAll(id, "")
	if err != nil {
		return
	}

	var aids []int

	// get article enid
	for _, p := range list.List {
		aids = append(aids, p.ID)
		if p.ClassID == id && p.ID == aid {
			aEnid = p.Enid
		}
	}
	if !utils.Contains(aids, aid) {
		err = errors.New("找不到该文章 ID，请检查输入是否正确")
		return
	}

	info, err = getService().ArticleInfo(aEnid, 1)
	if err != nil {
		return
	}
	return
}

// ArticleDetail article detail
func ArticleDetail(id, aid int) (detail *services.ArticleDetail, enId string, err error) {
	info, aEnid, err := ArticleInfo(id, aid)
	if err != nil {
		return
	}
	enId = aEnid
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, aEnid, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}

// ArticleCommentList article comment list
func ArticleCommentList(enId, sort string, page, limit int) (list *services.ArticleCommentList, err error) {
	list, err = getService().ArticleCommentList(enId, sort, page, limit)
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
