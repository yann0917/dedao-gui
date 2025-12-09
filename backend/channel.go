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

// ChannelTopicDetail 根据主题ID获取主题详情列表
func (a *App) ChannelTopicDetail(productID int) (topic *services.ChannelTopicCategory, err error) {
	data, err := app.ChannelTopicDetail(productID)
	if err != nil {
		return
	}
	topic = data.Topic
	return
}
