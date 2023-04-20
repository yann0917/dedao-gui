package backend

import "github.com/yann0917/dedao-gui/backend/services"

func (a *App) GetHomeInitialState() (state services.HomeInitState, err error) {
	state, err = Instance.GetHomeInitialState()
	return
}

func (a *App) SearchHot() (list *services.SearchTot, err error) {
	list, err = Instance.SearchHot()
	return
}

// SunflowerLabelList 首页导航标签列表
// 2-好看又好查的电子书, 4-精选课程
func (a *App) SunflowerLabelList(nType int) (list *services.SunflowerLabelList, err error) {
	list, err = Instance.SunflowerLabelList(nType)
	return
}

func (a *App) SunflowerLabelContent(enID string, nType, page, pageSize int) (list *services.SunflowerContent, err error) {
	list, err = Instance.SunflowerLabelContent(enID, nType, page, pageSize)
	return
}

func (a *App) SunflowerResourceList() (list *services.SunflowerResourceList, err error) {
	list, err = Instance.SunflowerResourceList()
	return
}

func (a *App) AlgoFilter(param services.AlgoFilterParam) (resp *services.AlgoFilterResp, err error) {
	resp, err = Instance.AlgoFilter(param)
	return
}

func (a *App) AlgoProduct(param services.AlgoFilterParam) (resp *services.AlgoProductResp, err error) {
	resp, err = Instance.AlgoProduct(param)
	return
}
