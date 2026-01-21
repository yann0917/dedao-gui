package app

import (
	"github.com/yann0917/dedao-gui/backend/services"
)

// CourseType 课程分类
func CourseType() (list *services.CourseCategoryList, err error) {
	list, err = getService().CourseType()
	return
}

// GetNavbar 获取导航配置
func GetNavbar() (data *services.NavbarData, err error) {
	data, err = getService().GetNavbar()
	return
}

func CourseList(category, order, filter string, page, limit int) (list *services.CourseList, err error) {
	list, err = getService().CourseList(category, order, filter, page, limit)
	if err != nil {
		return
	}
	return
}

func CourseGroupList(category, order, filter string, groupID, page, limit int) (list *services.CourseList, err error) {
	list, err = getService().CourseGroupList(category, order, filter, groupID, page, limit)
	if err != nil {
		return
	}
	return
}

// CourseDetail 已购课程详情
func CourseDetail(category string, id int) (detail *services.Course, err error) {
	detail, err = getService().CourseDetail(category, id)
	if err != nil {
		return
	}
	return
}

// CourseInfoByEnid 已购课程列表
func CourseInfoByEnid(enID string) (info *services.CourseInfo, err error) {
	info, err = getService().CourseInfo(enID)
	if err != nil {
		return
	}
	return
}

// OutsideDetail 名家讲书课程详情
func OutsideDetail(enid string) (detail *services.OutsideDetail, err error) {
	detail, err = getService().OutsideDetail(enid)
	if err != nil {
		return
	}
	return
}
