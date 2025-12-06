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

// ChannelCategoryList 获取学习圈分类列表（根据分类ID获取主题内容列表）
func ChannelCategoryList(channelID int, categoryID int) (topics []services.ChannelTopicCategory, err error) {
	homepage, err := getService().ChannelHomepage(channelID)
	if err != nil {
		return
	}

	// 查找指定分类
	for _, category := range homepage {
		if category.CategoryID == categoryID {
			topics = category.List
			return topics, nil
		}
	}

	return topics, nil
}
