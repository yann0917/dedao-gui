package services

type LiveTabList struct {
	List []LiveTab `json:"list"`
}

type LiveTab struct {
	TabName  string `json:"tab_name"`
	TabCount int    `json:"tab_count"`
	LiveType int    `json:"live_type"`
}

type LiveList struct {
	IsMore int    `json:"is_more"`
	List   []Live `json:"list"`
}

type Live struct {
	Status               int    `json:"status"`
	OnlineNum            int    `json:"online_num"`
	ReservationNum       int    `json:"reservation_num"`
	Title                string `json:"title"`
	RoomId               int    `json:"room_id"`
	Intro                string `json:"intro"`
	Starttime            int    `json:"starttime"`
	StarttimeDesc        string `json:"starttime_desc"`
	Duration             int    `json:"duration"`
	Endtime              int    `json:"endtime"`
	Id                   int    `json:"id"`
	Type                 int    `json:"type"`
	LogId                int    `json:"log_id"`
	LogType              string `json:"log_type"`
	LiveType             string `json:"live_type"`
	LivePrivilegeTips    string `json:"live_privilege_tips"`
	HomeImg              string `json:"home_img"`
	Author               string `json:"author"`
	PvNum                int    `json:"pv_num"`
	PrivilegeStatus      int    `json:"privilege_status"`
	AliasId              string `json:"alias_id"`
	AppointmentStatus    int    `json:"appointment_status"`
	PlaybackStatus       int    `json:"playback_status"`
	PublishStatus        int    `json:"publish_status"`
	LiveDurationText     string `json:"live_duration_text"`
	VideoDuration        string `json:"video_duration"`
	LivePrivilegeType    int    `json:"live_privilege_type"`
	Logo                 string `json:"logo"`
	TimeReport           int    `json:"time_report"`
	HideOnlineNumber     int    `json:"hide_online_number"`
	LivePv               int    `json:"live_pv"`
	LiveViewers          int    `json:"live_viewers"`
	ShowPv               int    `json:"show_pv"`
	IsPrivilegeLive      bool   `json:"is_privilege_live"`
	InviteCount          int    `json:"invite_count"`
	CanWatch             bool   `json:"can_watch"`
	HasBuy               bool   `json:"has_buy"`
	PrivilegeLiveTag     string `json:"privilege_live_tag"`
	SubscribeSummary     string `json:"subscribe_summary"`
	DdUrl                string `json:"dd_url"`
	ShowSubscribeSummary int    `json:"show_subscribe_summary"`
	AlertTips            string `json:"alert_tips"`
	Currenttime          int    `json:"currenttime"`
	VideoCoverM3U8       string `json:"video_cover_m3u8"`
	LiveCoverM3U8        string `json:"live_cover_m3u8"`
	WebPcMediaToken      string `json:"web_pc_media_token"`
	LdFlv                string `json:"ld_flv"`
	VideoCoverMediaId    int64  `json:"video_cover_media_id"`
	LastStartTimestamp   int    `json:"last_start_timestamp"`
	LastEndTimestamp     int    `json:"last_end_timestamp"`

	PrivilegeProduct LivePrivilegeProduct `json:"privilege_product"`
}

type LivePrivilegeProduct struct {
	Id              int    `json:"id"`
	ProductId       int    `json:"product_id"`
	ProductSubId    int    `json:"product_sub_id"`
	ProductType     int    `json:"product_type"`
	Title           string `json:"title"`
	Ddurl           string `json:"ddurl"`
	LiveStartTime   int    `json:"live_start_time"`
	LiveId          int    `json:"live_id"`
	HasPrivilege    bool   `json:"has_privilege"`
	AliasName       string `json:"alias_name"`
	AliasId         string `json:"alias_id"`
	IsPrivilegeLive bool   `json:"is_privilege_live"`
	ProductGroupId  int    `json:"product_group_id"`
}
type LiveCheck struct {
	Status             int    `json:"status"`
	Token              string `json:"token"`
	LivePrivilegeTips  string `json:"live_privilege_tips"`
	LiveType           string `json:"live_type"`
	Logo               string `json:"logo"`
	ExpiredSeconds     int    `json:"expired_seconds"`
	ProductType        int    `json:"product_type"`
	ProductId          int    `json:"product_id"`
	Ddurl              string `json:"ddurl"`
	NeedBuy            bool   `json:"need_buy"`
	ButtonText         string `json:"button_text"`
	ProductSubid       int    `json:"product_subid"`
	LiveProductBuyTips string `json:"live_product_buy_tips"`
	InviteCount        int    `json:"invite_count"`
	IsPrivilegeLive    bool   `json:"is_privilege_live"`
	ErrorMsg           string `json:"error_msg"`
	HasBuy             bool   `json:"has_buy"`
}

type LiveBase struct {
	LiveLectureList     []LiveLecture    `json:"live_lecture_list"`
	LiveOutline         LiveOutline      `json:"live_outline"`
	RoomDetailDdurl     string           `json:"room_detail_ddurl"`
	Status              int              `json:"status"`
	ShareImageSquare    string           `json:"share_image_square"`
	Logo                string           `json:"logo"`
	PosterImage         string           `json:"poster_image"`
	Intro               string           `json:"intro"`
	Title               string           `json:"title"`
	Starttimestamp      int              `json:"starttimestamp"`
	ShareTitle          string           `json:"share_title"`
	ShareSummary        string           `json:"share_summary"`
	Id                  int              `json:"id"`
	ShareUrl            string           `json:"share_url"`
	Booking             LiveBooking      `json:"booking"`
	LiveSeriesIds       string           `json:"live_series_ids"`
	PlaybackStatus      int              `json:"playback_status"`
	PublishStatus       int              `json:"publish_status"`
	LiveSeriresInfo     interface{}      `json:"live_serires_info"`
	LiveType            string           `json:"live_type"`
	Ptype               int              `json:"ptype"`
	LastStartTime       string           `json:"last_start_time"`
	LastEndTime         string           `json:"last_end_time"`
	MaxNum              int              `json:"max_num"`
	Uv                  int              `json:"uv"`
	LivePv              int              `json:"live_pv"`
	LiveViewers         int              `json:"live_viewers"`
	ShowPv              int              `json:"show_pv"`
	InviteCount         int              `json:"invite_count"`
	IsPrivilegeLive     bool             `json:"is_privilege_live"`
	AuthorText          string           `json:"author_text"`
	LiveArticleInfo     LiveArticleInfo  `json:"live_article_info"`
	LiveActivityInfo    LiveActivityInfo `json:"live_activity_info"`
	IsActivityLive      bool             `json:"is_activity_live"`
	YouzanVipProductUrl string           `json:"youzan_vip_product_url"`
	LiveUv              int              `json:"live_uv"`
	LiveMedalModule     interface{}      `json:"live_medal_module"`
	LivePrivilegeType   int              `json:"live_privilege_type"`
	LiveRecordNum       string           `json:"live_record_num"`
}

type LiveLecture struct {
	TeacherUid      int    `json:"teacher_uid"`
	TeacherName     string `json:"teacher_name"`
	TeacherIntro    string `json:"teacher_intro"`
	Avatar          string `json:"avatar"`
	ProfileUrl      string `json:"profile_url"`
	IsBigv          int    `json:"is_bigv"`
	FllowStatus     int    `json:"fllow_status"`
	FollowCount     int    `json:"follow_count"`
	FansCount       int    `json:"fans_count"`
	TeacherUidAlias string `json:"teacher_uid_alias"`
}

type LiveOutline struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type LiveBooking struct {
	BookStatus int `json:"bookStatus"`
	TotalNum   int `json:"totalNum"`
}

type LiveArticleInfo struct {
	ArticleId          int    `json:"article_id"`
	ArticleDetailDdurl string `json:"article_detail_ddurl"`
}

type LiveActivityInfo struct {
	LiveActivityDdurl    string `json:"live_activity_ddurl"`
	LiveActivityType     int    `json:"live_activity_type"`
	LivePlaybackDisabled int    `json:"live_playback_disabled"`
	LiveActivityId       string `json:"live_activity_id"`
}

func (s *Service) LiveTabList() (list *LiveTabList, err error) {
	body, err := s.reqLiveTabList()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}

	return
}

func (s *Service) LiveList(liveType, page, limit int) (list *LiveList, err error) {
	body, err := s.reqLiveList(liveType, page, limit)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

func (s *Service) LiveCheck(aliasID, inviteCode string) (detail *LiveCheck, err error) {
	body, err := s.reqLiveCheck(aliasID, inviteCode)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

func (s *Service) LiveBase(aliasID string) (detail *LiveBase, err error) {
	body, err := s.reqLiveBase(aliasID)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}
