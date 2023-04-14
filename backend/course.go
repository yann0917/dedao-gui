package backend

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/yann0917/dedao/app"
	"github.com/yann0917/dedao/services"
)

func (a *App) CourseCategory() (list []services.CourseCategory, err error) {
	result, err := app.CourseType()
	if err != nil {
		return
	}
	list = result.Data.List
	return
}

func (a *App) CourseList(category, order string, page, limit int) (list *services.CourseList, err error) {
	list, err = app.CourseList(category, order, page, limit)
	if err != nil {
		return
	}

	return
}

func (a *App) CourseInfo(enid string) (info *services.CourseInfo, err error) {
	info, err = app.CourseInfoByEnid(enid)
	return
}

func (a *App) ArticleList(enid, chapterID string, count, maxID int) (list *services.ArticleList, err error) {
	list, err = app.ArticleList(enid, chapterID, count, maxID)
	if err != nil {
		return
	}

	return
}

// ArticleDetail
// enid article enid  or odob audioAliasID, aType 1-course article, 2-odob article
func (a *App) ArticleDetail(aType int, aEnid string) (markdown string, err error) {
	detail, err := app.ArticleDetailByEnid(aType, aEnid)
	if err != nil {
		return
	}

	var content []services.Content
	err = jsoniter.UnmarshalFromString(detail.Content, &content)
	if err != nil {
		return
	}
	markdown = app.ContentsToMarkdown(content)
	return
}

func (a *App) GetVolcPlayAuthToken(mediaID, securityToken string) (info *services.MediaVolc, err error) {
	info, err = Instance.GetVolcPlayAuthToken(mediaID, securityToken)
	// fmt.Println(info)
	// fmt.Println(err)
	return
}

func (a *App) GetVolcPlayInfo(query string) (info *services.VodPlayInfoResp, err error) {
	info, err = Instance.GetVolcPlayInfo(query)
	if err != nil {
		return
	}
	return
}
