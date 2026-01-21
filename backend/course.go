package backend

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
)

func (a *App) CourseCategory() (list []services.CourseCategory, err error) {
	result, err := app.CourseType()
	if err != nil {
		return
	}
	list = result.Data.List
	return
}

func (a *App) GetNavbar() (data *services.NavbarData, err error) {
	data, err = app.GetNavbar()
	return
}

func (a *App) CourseList(category, order, filter string, page, limit int) (list *services.CourseList, err error) {
	list, err = app.CourseList(category, order, filter, page, limit)
	if err != nil {
		return
	}

	return
}

func (a *App) CourseGroupList(category, order, filter string, groupID, page, limit int) (list *services.CourseList, err error) {
	list, err = app.CourseGroupList(category, order, filter, groupID, page, limit)
	if err != nil {
		return
	}
	return
}

func (a *App) CourseInfo(enid string) (info *services.CourseInfo, err error) {
	info, err = app.CourseInfoByEnid(enid)
	return
}

func (a *App) OutsideDetail(enid string) (detail *services.OutsideDetail, err error) {
	detail, err = app.OutsideDetail(enid)
	if err != nil {
		return
	}
	return
}

func (a *App) ArticleList(enid, chapterID string, count, maxID int, reverse bool) (list *services.ArticleList, err error) {
	list, err = app.ArticleList(enid, chapterID, count, maxID, reverse)
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

func (a *App) GetArticleIntro(aType int, enid string) (intro *services.ArticleIntro, err error) {
	info, err := Instance.ArticleInfo(enid, aType)
	if err != nil {
		return
	}
	intro = &info.ArticleInfo
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
