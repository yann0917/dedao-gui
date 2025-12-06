package backend

import (
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
)

// ChannelInfo 获取学习圈频道信息
func (a *App) ChannelInfo(channelID int) (info *services.ChannelInfo, err error) {
	return app.ChannelInfo(channelID)
}

// ChannelHomepage 获取学习圈频道首页分类
func (a *App) ChannelHomepage(channelID int) (cats []services.ChannelHomepageCategory, err error) {
	return app.ChannelHomepage(channelID)
}

// ChannelVipInfo 获取学习圈VIP/权限信息
func (a *App) ChannelVipInfo(channelID int) (info *services.ChannelVipInfo, err error) {
	return app.ChannelVipInfo(channelID)
}

// ChannelCategoryList 获取学习圈分类列表（根据分类ID获取主题内容列表）
func (a *App) ChannelCategoryList(channelID int, categoryID int) (topics []services.ChannelTopicCategory, err error) {
	return app.ChannelCategoryList(channelID, categoryID)
}
