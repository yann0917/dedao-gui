package app

import (
	"github.com/yann0917/dedao-gui/backend/services"
)

// AudioDetail 听书音频简介
func AudioDetail(id string) (detail *services.AudioInfoResp, err error) {
	detail, err = getService().AudioDetail(id)
	if err != nil {
		return
	}
	return
}

func AudioDetailAlias(aliasID string) (detail *services.Audio, err error) {
	detail, err = getService().AudioDetailAlias(aliasID)
	if err != nil {
		return
	}
	return
}

// OdobShelfAdd 听书加入书架
func OdobShelfAdd(enIDs []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = getService().OdobShelfAdd(enIDs)
	return
}
