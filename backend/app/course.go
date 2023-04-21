package app

import (
	"github.com/yann0917/dedao-gui/backend/services"
)

// CourseType 课程分类
func CourseType() (list *services.CourseCategoryList, err error) {
	list, err = getService().CourseType()
	return
}

func CourseList(category, order string, page, limit int) (list *services.CourseList, err error) {
	list, err = getService().CourseList(category, order, page, limit)
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
