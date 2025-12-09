package app

import (
	"github.com/yann0917/dedao-gui/backend/services"
)

// ChannelInfo 获取学习圈频道信息
func ChannelInfo(channelID int) (info *services.ChannelInfo, err error) {
	info, err = getService().ChannelInfo(channelID)
	if err != nil {
		return
	}
	return
}

// ChannelHomepage 获取学习圈频道首页分类
func ChannelHomepage(channelID int) (cats []services.ChannelHomepageCategory, err error) {
	cats, err = getService().ChannelHomepage(channelID)
	if err != nil {
		return
	}
	return
}

// ChannelVipInfo 获取学习圈VIP/权限信息
func ChannelVipInfo(channelID int) (info *services.ChannelVipInfo, err error) {
	info, err = getService().ChannelVipInfo(channelID)
	if err != nil {
		return
	}
	return
}

// ChannelTopicDetail 根据主题ID获取主题详情
func ChannelTopicDetail(productID int) (topic *services.ChannelTopicDetail, err error) {
	topic, err = getService().ChannelTopicDetail(productID)
	if err != nil {
		return
	}
	return
}
