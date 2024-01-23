package services

type AudioInfo struct {
	ID             int            `json:"id"`
	AudioID        string         `json:"audio_id"`
	Title          string         `json:"title"`
	Icon           string         `json:"icon"`
	Duration       int            `json:"duration"`
	AudioPrice     string         `json:"audio_price"`
	AudioSummary   string         `json:"audio_summary"`
	PublishTime    int            `json:"publish_time"`
	IsVipGived     bool           `json:"is_vip_gived"`
	IsVip          bool           `json:"is_vip"`
	AgencyDetail   AgencyDetail   `json:"agency_detail"`
	LogID          int            `json:"log_id"`
	IsBuy          bool           `json:"is_buy"`
	InBookrack     bool           `json:"in_bookrack"`
	Progress       int            `json:"progress"`
	Tag            []string       `json:"tag"`
	TopicSummary   []TopicSummary `json:"topic_summary"`
	IsLimitFree    bool           `json:"is_limit_free"`
	TopicEncodeId  string         `json:"topic_encode_id"`
	H5ShareUrl     string         `json:"h5_share_url"`
	CanShare       bool           `json:"can_share"`
	LimitFree      LimitFree      `json:"limit_free"`
	VipEndTime     int            `json:"vip_end_time"`
	ButtonType     int            `json:"button_type"`
	HasPlayAuth    bool           `json:"has_play_auth"`
	Rank           Rank           `json:"rank"`
	LearnCount     int            `json:"learn_count"`
	LearnCountDesc string         `json:"learn_count_desc"`
}

type AgencyDetail struct {
	ID              int    `json:"id"`
	IDStr           string `json:"id_str"`
	Name            string `json:"name"`
	Intro           string `json:"intro"`
	QcgID           int    `json:"qcg_id"`
	QcgMemberName   string `json:"qcg_member_name"`
	QcgMemberAvatar string `json:"qcg_member_avatar"`
	BookCount       int    `json:"book_count"`
	Uv              int    `json:"uv"`
	Status          int    `json:"status"`
}

type TopicSummary struct {
	ID        int    `json:"id"`
	TopicID   int    `json:"topic_id"`
	Title     string `json:"title"`
	SubTitle  string `json:"sub_title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LimitFree struct {
	LimitFreeExpire bool `json:"limit_free_expire"`
	IsLimitFree     bool `json:"is_limit_free"`
	IsRedPacket     bool `json:"is_red_packet"`
	FreeBeginTime   int  `json:"free_begin_time"`
	FreeEndTime     int  `json:"free_end_time"`
	FreeMaximum     int  `json:"free_maximum"`
	ConsumeNum      int  `json:"consume_num"`
	VipEndTime      int  `json:"vip_end_time"`
}

type Rank struct {
	RankNumber int    `json:"rank_number"`
	RankName   string `json:"rank_name"`
	RankType   int    `json:"rank_type"`
	RankDesc   string `json:"rank_desc"`
	DdUrl      string `json:"dd_url"`
}

type Quality struct {
	Desc    string        `json:"desc"`
	List    []QualityUser `json:"list"`
	Paytime string        `json:"paytime"`
}
type QualityUser struct {
	Avatar string `json:"avatar"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
}
type AudioInfoResp struct {
	AudioInfo AudioInfo `json:"detail"`
	Quality   Quality   `json:"quality"`
}

// AudioDetail get odob introduce
func (s *Service) AudioDetail(ID string) (detail *AudioInfoResp, err error) {
	body, err := s.reqAudioDetail(ID)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

// OdobShelfAdd add ebook shelf
func (s *Service) OdobShelfAdd(ids []string) (resp *EbookShelfAddResp, err error) {

	body, err := s.reqOdobShelfAdd(ids)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}
