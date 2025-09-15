package services

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// reqGetLoginAccessToken 扫码请求token
func (s *Service) reqGetLoginAccessToken(csrfToken string) (string, error) {
	resp, err := s.client.R().
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeaderVerbatim("Xi-Csrf-Token", csrfToken).
		Post("/loginapi/getAccessToken")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return "", err
	}
	accessToken := resp.String()
	return accessToken, err
}

// reqGetQrcode 扫码登录二维码
// token: X-Oauth-Access-Token from /loginapi/getAccessToken
func (s *Service) reqGetQrcode(token string) (qr *QrCodeResp, err error) {
	_, err = s.client.R().
		SetHeaderVerbatim("X-Oauth-Access-Token", token).
		SetResult(&qr).
		Get("/oauth/api/embedded/qrcode")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return
	}
	return
}

// reqCheckLogin 轮询扫码登录结果
// token: X-Oauth-Access-Token from /loginapi/getAccessToken
// qrCode: qrCodeString from /oauth/api/embedded/qrcode
func (s *Service) reqCheckLogin(token, qrCode string) (check *CheckLoginResp, cookie string, err error) {
	resp, err := s.client.R().
		SetHeaderVerbatim("X-Oauth-Access-Token", token).
		SetBody(map[string]interface{}{
			"keepLogin": true,
			"pname":     "mobilesms",
			"qrCode":    qrCode,
			"scene":     "login",
		}).
		SetResult(&check).
		Post("/oauth/api/embedded/qrcode/check_login")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return
	}
	cookies := resp.Header().Values("Set-Cookie")
	cookie = strings.Join(cookies, "; ")
	return
}

// reqToken 请求token
func (s *Service) reqToken() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Get("/ddph/v2/token/create")
	return handleHTTPResponse(resp, err)
}

// reqUser 请求用户信息
func (s *Service) reqUser() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Get("/api/pc/user/info")

	return handleHTTPResponse(resp, err)
}

// reqCourseType 请求首页课程分类列表
func (s *Service) reqCourseType() (io.ReadCloser, error) {
	resp, err := s.client.R().Post("/api/hades/v1/index/detail")
	return handleHTTPResponse(resp, err)
}

// reqCourseList 请求课程列表
func (s *Service) reqCourseList(category, order string, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.R().SetBody(map[string]interface{}{
		"category":        category,
		"display_group":   true,
		"filter":          "all",
		"group_id":        0,
		"order":           order,
		"filter_complete": 0,
		"page":            page,
		"page_size":       limit,
		"sort_type":       "desc",
	}).Post("/api/hades/v2/product/list")
	return handleHTTPResponse(resp, err)
}

// reqOutsideDetail 请求名家讲书课程详情
func (s *Service) reqOutsideDetail(enid string) (io.ReadCloser, error) {
	resp, err := s.client.R().SetBody(map[string]interface{}{
		"product_enid": enid,
		"product_type": 1013,
	}).Post("pc/sunflower/v1/depot/outside/detail")
	return handleHTTPResponse(resp, err)
}

// reqCourseInfo 请求课程介绍
func (s *Service) reqCourseInfo(ID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"detail_id": ID,
			"is_login":  1,
		}).
		Post("/pc/bauhinia/pc/class/info")
	return handleHTTPResponse(resp, err)
}

// reqArticleList 请求文章列表
// chapterID = "" 获取所有的文章列表，否则只获取该章节的文章列表
func (s *Service) reqArticleList(ID, chapterID string, count, maxID int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"chapter_id":      chapterID,
			"count":           count,
			"detail_id":       ID,
			"include_edge":    false,
			"is_unlearn":      false,
			"max_id":          maxID,
			"max_order_num":   0,
			"reverse":         false,
			"since_id":        0,
			"since_order_num": 0,
			"unlearn_switch":  false,
		}).Post("/api/pc/bauhinia/pc/class/purchase/article_list")
	return handleHTTPResponse(resp, err)
}

// reqArticleCommentList 请求文章热门留言列表
// enId 文章 ID
// sort like-最热 create-最新
// source_type 65-课程文章 create-最新
func (s *Service) reqArticleCommentList(enId, sort string, page, limit, sType int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"detail_enid":  enId,
			"note_type":    2,
			"only_replied": false,
			"page":         page,
			"page_count":   limit,
			"sort_by":      sort,
			"source_type":  sType,
		}).Post("/pc/ledgers/notes/article_comment_list")
	return handleHTTPResponse(resp, err)
}

// reqArticleInfo 请求文章 token
// id article id or odob audioAliasID
func (s *Service) reqArticleInfo(ID string, aType int) (io.ReadCloser, error) {
	param := make(map[string]string)
	switch aType {
	case 1:
		param["detail_id"] = ID
	case 2:
		param["audio_alias_id"] = ID
	}
	resp, err := s.client.R().
		SetBody(param).Post("/pc/bauhinia/pc/article/info")
	return handleHTTPResponse(resp, err)
}

// reqArticleDetail 请求文章详情
func (s *Service) reqArticleDetail(token, appID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"token":  token,
			"appid":  appID,
			"is_new": "1",
		}).
		Get("/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqArticlePoint 请求文章重点
func (s *Service) reqArticlePoint(enid string, pType string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"article_id_hazy": enid,
			"product_type":    pType,
		}).Get("/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqAudioByAlias 请求音频详情
func (s *Service) reqAudioByAlias(ids string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"ids":            ids,
			"get_extra_data": 1,
		}).
		Post("/pc/bauhinia/v1/audio/mutiget_by_alias")
	return handleHTTPResponse(resp, err)
}

// reqAudioDetail 每天听本书书简介
func (s *Service) reqAudioDetail(topicID string) (io.ReadCloser, error) {
	resp, err := s.client.R().SetQueryParams(map[string]string{
		"topic_id_str": topicID,
	}).
		Get("pc/odob/pc/audio/detail")
	return handleHTTPResponse(resp, err)
}

// reqOdobAudioDetail 请求每天听本书书 音频 info
func (s *Service) reqOdobAudioDetail(aliasID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"alias_id": aliasID,
		}).
		Post("pc/odob/pc/audio/detail/alias")

	return handleHTTPResponse(resp, err)
}

// reqOdobShelfAdd 请求听书加入书架
func (s *Service) reqOdobShelfAdd(enIds []string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"book_enids": enIds,
		}).
		Post("pc/odob/v2/bookrack/pc/add")

	return handleHTTPResponse(resp, err)
}

// reqEbookDetail 请求电子书详情
func (s *Service) reqEbookDetail(enid string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParam("id", enid).
		Get("/pc/ebook2/v1/pc/detail")

	return handleHTTPResponse(resp, err)
}

// reqEbookReadToken 请求电子书阅读 token
func (s *Service) reqEbookReadToken(enid string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]string{
			"id": enid,
		}).
		Post("/api/pc/ebook2/v1/pc/read/token")
	return handleHTTPResponse(resp, err)
}

// reqEbookInfo 请求电子书 info
func (s *Service) reqEbookInfo(token string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParam("token", token).
		Get("/ebk_web/v1/get_book_info")
	return handleHTTPResponse(resp, err)
}

// reqEbookPages 获取页面详情
func (s *Service) reqEbookPages(chapterID, token string, index, count, offset int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"chapter_id":  chapterID,
			"count":       count,
			"index":       index,
			"offset":      offset,
			"orientation": 0,
			"config": map[string]interface{}{
				"density":         1,
				"direction":       0,
				"font_name":       "pingfang",
				"font_scale":      1,
				"font_size":       16,
				"height":          200000,
				"line_height":     "2em",
				"margin_bottom":   20,
				"margin_left":     20,
				"margin_right":    20,
				"margin_top":      0,
				"paragraph_space": "1em",
				"platform":        1,
				"width":           60000,
			},
			"token": token,
		}).
		Post("/ebk_web_go/v2/get_pages")
	return handleHTTPResponse(resp, err)
}

// reqEbookInfo 请求电子书vip info
func (s *Service) reqEbookVIPInfo() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Post("/api/pc/ebook2/v1/vip/info")
	return handleHTTPResponse(resp, err)
}

// reqEbookCommentList 请求电子评分及书评列表
// sort like_count
func (s *Service) reqEbookCommentList(enid, sort string, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"enid":      enid,
			"page":      page,
			"page_size": limit,
			"sort":      sort,
			"ptype":     2,
		}).
		Post("/pc/ebook2/v1/comment/list")
	return handleHTTPResponse(resp, err)
}

// reqEbookDetail 请求电子书加入书架
func (s *Service) reqEbookShelfAdd(enIds []string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"book_enids": enIds,
		}).
		Post("/api/pc/ebook2/v1/bookshelf/add")

	return handleHTTPResponse(resp, err)
}

// reqEbookDetail 请求电子书移出书架
func (s *Service) reqEbookShelfRemove(ids []string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"pids":  ids,
			"ptype": 2,
		}).
		Post("/api/pc/hades/v1/product/remove")

	return handleHTTPResponse(resp, err)
}

// reqOdobVIPInfo 请求每天听本书书 vip info
func (s *Service) reqOdobVIPInfo() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Post("pc/odob/v2/vipuser/vip_card_info")

	return handleHTTPResponse(resp, err)
}

// reqTopicAll 请求推荐话题列表
func (s *Service) reqTopicAll(page, limit int) (io.ReadCloser, error) {
	r := s.client.R()
	r = r.SetBody(map[string]int{
		"page_id": page,
		"limit":   limit,
	})
	resp, err := r.Post("/pc/ledgers/topic/all")
	return handleHTTPResponse(resp, err)
}

// reqTopicDetail 请求话题详情
func (s *Service) reqTopicDetail(topicID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"incr_view_count": true,
			"topic_id_hazy":   topicID,
		}).Post("/pc/ledgers/topic/detail")
	return handleHTTPResponse(resp, err)
}

// reqTopicNotesList 请求话题笔记列表
func (s *Service) reqTopicNotesList(topicID string, isElected bool, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"count":         limit,
			"is_elected":    isElected,
			"page_id":       page,
			"version":       2,
			"topic_id_hazy": topicID,
		}).Post("/pc/ledgers/topic/notes/list")
	return handleHTTPResponse(resp, err)
}

// reqTopicNotesTimeline 请求话题笔记时间线列表
func (s *Service) reqTopicNotesTimeline(maxID string) (io.ReadCloser, error) {
	param := make(map[string]interface{})
	if maxID != "" {
		param["max_id"] = maxID
	}
	param["count"] = 20
	param["load_chain"] = true
	param["version"] = 2

	resp, err := s.client.R().
		SetBody(param).Post("/pc/ledgers/notes/friends_timeline")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqLiveTabList() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Post("/api/pc/ddlive/v2/pc/home/live/tablist")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqLiveList(liveType, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"live_type": liveType,
			"page":      page,
			"page_size": limit,
		}).Post("/api/pc/ddlive/v2/pc/home/live/list")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqLiveCheck(aliasID, inviteCode string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"alias_id":    aliasID,
			"invite_code": inviteCode,
		}).Post("/api/pc/ddlive/v2/pc/live/check")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqLiveBase(aliasID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"alias_id": aliasID,
		}).Post("/pc/ddlive/v2/pc/live/base")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqTimeReport(data interface{}) (io.ReadCloser, error) {
	// data
	// [{
	// 	"resource": {
	// 		"product_type": 8888
	// 	},
	// 	"time": {
	// 		"model": 1,
	// 		"speed": 1,
	// 		"scope": "1666832665-1666832669,1666832726-1666832785",
	// 		"extra": "active"
	// 		}
	// }]
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"data":      data,
			"os":        "H5",
			"platform":  6,
			"send_time": time.Now().UnixMilli(),
		}).Post("/prime/v1/producer/time/report")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqVolc(mediaID, securityToken string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"media_id_str":   mediaID,
			"security_token": securityToken,
		}).Post("/media_gate/gate/api/v1/volc")
	return handleHTTPResponse(resp, err)
}

func (s *Service) reqVolcGetPlayInfo(query string) (io.ReadCloser, error) {
	cookies := s.client.Cookies
	client := resty.New()
	client.SetCookies(cookies).
		SetHeaderVerbatim("User-Agent", UserAgent).
		SetHeaderVerbatim("Xi-DT", "web")
	client.SetBaseURL("https://vod.volcengineapi.com")
	resp, err := client.
		R().
		SetQueryString("Action=GetPlayInfo&Version=2020-08-01&" + query).
		Get("")
	return handleHTTPResponse(resp, err)
}

// reqSearchHot search hot
func (s *Service) reqSearchHot() (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"is_login": 0,
		}).
		Post("/api/search/pc/hot")
	return handleHTTPResponse(resp, err)
}

// reqSunflowerLabelList index label list
// nType 2-好看又好查的电子书, 4-精选课程
func (s *Service) reqSunflowerLabelList(nType int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"label_count": -1,
			"nav_type":    nType,
		}).
		Post("/pc/sunflower/v1/label/list")
	return handleHTTPResponse(resp, err)
}

// reqSunflowerLabelContent index label contents
// enID /pc/sunflower/v1/label/list 接口 enid
// nType 2-好看又好查的电子书, 4-精选课程
func (s *Service) reqSunflowerLabelContent(enID string, nType, page, pageSize int) (io.ReadCloser, error) {
	var resType string
	if nType == 2 {
		resType = "2"
	} else if nType == 4 {
		resType = "66"
	}
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"enid":        enID,
			"nav_type":    nType,
			"page":        page,
			"page_size":   pageSize,
			"result_type": resType,
		}).
		Post("/pc/sunflower/v1/label/content")
	return handleHTTPResponse(resp, err)
}

// reqSunflowerResourceList 免费专区
func (s *Service) reqSunflowerResourceList() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Get("/pc/sunflower/v1/resource/list")
	return handleHTTPResponse(resp, err)
}

// reqAlgoFilter index filter
// 1. 获取全部，参数：cName default "心理学"， pageSize = 18, pType="66"
func (s *Service) reqAlgoFilter(param AlgoFilterParam) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(param).
		Post("/pc/label/v2/algo/pc/filter/list")
	return handleHTTPResponse(resp, err)
}

// reqAlgoProduct index label contents
func (s *Service) reqAlgoProduct(param AlgoFilterParam) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBody(param).
		Post("/pc/label/v2/algo/pc/product/list")
	return handleHTTPResponse(resp, err)
}
