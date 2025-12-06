package services

// ChannelInfo 学习圈频道信息
type ChannelInfo struct {
	ChannelID            int               `json:"channel_id"`
	IsFree               bool              `json:"is_free"`
	IsSharable           bool              `json:"is_sharable"`
	IsShowEquity         bool              `json:"is_show_equity"`
	SubscriptionPType    int               `json:"subscription_ptype"`
	Theme                string            `json:"theme"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	Logo                 string            `json:"logo"`
	Badge                string            `json:"badge"`
	ThemeColor           string            `json:"theme_color"`
	Host                 ChannelPerson     `json:"host"`
	Hosts                []ChannelPerson   `json:"hosts"`
	Guests               []ChannelPerson   `json:"guests"`
	EquityTitle          string            `json:"equity_title"`
	EquityDescription    string            `json:"equity_description"`
	GuestsURL            string            `json:"guests_url"`
	EquityURL            string            `json:"equity_url"`
	PurchaseURL          string            `json:"purchase_url"`
	GuestsDDURL          string            `json:"guests_ddurl"`
	EquityDDURL          string            `json:"equity_ddurl"`
	PurchaseDDURL        string            `json:"purchase_ddurl"`
	Statistics           ChannelStatistics `json:"statistics"`
	MainButton           string            `json:"main_button"`
	ViceButton           string            `json:"vice_button"`
	ViceLine             string            `json:"vice_line"`
	StyleSheet           ChannelStyleSheet `json:"style_sheet"`
	EquityManageDDURL    string            `json:"equity_manage_ddurl"`
	IsVIP                bool              `json:"is_vip"`
	IsPop                bool              `json:"is_pop"`
	ShowAIAssistant      bool              `json:"show_ai_assistant"`
	ShowNewClassesLayout bool              `json:"show_new_classes_layout"`
	ShowCouponRenewal    bool              `json:"show_coupon_renewal"`
	CountdownTips        string            `json:"countdown_tips"`
}

// ChannelPerson 频道人物（主持人/嘉宾通用）
type ChannelPerson struct {
	UID    int    `json:"uid"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
	VStat  int    `json:"v_stat"`
}

// ChannelStatistics 学习圈数据统计
type ChannelStatistics struct {
	TotalSubscribers int                   `json:"total_subscribers"`
	SellingPoints    []ChannelSellingPoint `json:"selling_points"`
	ContentQuantity  int                   `json:"content_quantity"`
	Messages         []string              `json:"messages"`
	Tips             []string              `json:"tips"`
}

// ChannelSellingPoint 卖点
type ChannelSellingPoint struct {
	CountDesc string `json:"count_desc"`
	Name      string `json:"name"`
}

// ChannelStyleSheet 展示样式
type ChannelStyleSheet struct {
	HostBG     string `json:"host_bg"`
	HostBorder string `json:"host_border"`
}

// ChannelHomepageCategory 频道首页分类
type ChannelHomepageCategory struct {
	CategoryID       int                    `json:"category_id"`
	CategoryName     string                 `json:"category_name"`
	CategoryIcon     string                 `json:"category_icon"`
	CategoryDarkIcon string                 `json:"category_dark_icon"`
	List             []ChannelTopicCategory `json:"list"`
}

// ChannelTopicCategory 分类下的主题
type ChannelTopicCategory struct {
	ID           int           `json:"id"`
	ChannelID    int           `json:"channel_id"`
	Title        string        `json:"title"`
	Icon         string        `json:"icon"`
	NightIcon    string        `json:"night_icon"`
	Intro        string        `json:"intro"`
	Description  string        `json:"description"`
	Status       int           `json:"status"`
	SortNo       int           `json:"sort_no"`
	ContentType  string        `json:"content_type"`
	ShowHomepage int           `json:"show_homepage"`
	DDURL        string        `json:"dd_url"`
	Items        []ChannelItem `json:"items"`
	Length       int           `json:"length"`
}

// ChannelItem 主题内条目（课程/视频/直播等）
type ChannelItem struct {
	ProductType        int               `json:"product_type"`
	ProductID          int               `json:"product_id"`
	EnID               string            `json:"en_id"`
	ClassID            int               `json:"class_id"`
	ClassEnID          string            `json:"class_en_id"`
	Title              string            `json:"title"`
	Cover              string            `json:"cover"`
	Logo               string            `json:"logo"`
	Summary            string            `json:"summary"`
	DifficultyLevel    int               `json:"difficulty_level"`
	DDURL              string            `json:"ddurl"`
	Status             string            `json:"status"`
	HasVideo           bool              `json:"has_video"`
	HasLive            bool              `json:"has_live"`
	Duration           int               `json:"duration"`
	Progress           float64           `json:"progress"`
	PublishTime        int               `json:"publish_time"`
	LearnCount         int               `json:"learn_count"`
	IsRead             bool              `json:"is_read"`
	LiveStatus         int               `json:"live_status"`
	LivePlaybackStatus int               `json:"live_playback_status"`
	LiveStatusTips     string            `json:"live_status_tips"`
	Source             ChannelItemSource `json:"source"`
}

// ChannelItemSource 条目来源信息
type ChannelItemSource struct {
	SourceName string `json:"source_name"`
	SourceType int    `json:"source_type"`
	SourceID   int    `json:"source_id"`
}

// ChannelVipInfo 学习圈VIP/权限信息
type ChannelVipInfo struct {
	UID               int            `json:"uid"`
	Nickname          string         `json:"nickname"`
	Slogan            string         `json:"slogan"`
	Avatar            string         `json:"avatar"`
	AvatarS           string         `json:"avatar_s"`
	VStateValue       int            `json:"v_state_value"`
	IsVIP             bool           `json:"is_vip"`
	BeginTime         int64          `json:"begin_time"`
	EndTime           int64          `json:"end_time"`
	ExpireTime        int64          `json:"expire_time"`
	CurrentTime       int64          `json:"current_time"`
	SurplusDays       int            `json:"surplus_days"`
	SubscribedDays    int            `json:"subscribed_days"`
	IsEverSubscribed  bool           `json:"is_ever_subscribed"`
	IsExpire          bool           `json:"is_expire"`
	CurrentWorkCardID int            `json:"current_work_card_id"`
	ErrTips           string         `json:"err_tips"`
	Medals            []ChannelMedal `json:"medals"`
	DefaultVipCard    DefaultVipCard `json:"default_vip_card"`
	EquityIntro       string         `json:"equity_intro"`
}

// ChannelMedal 勋章
type ChannelMedal struct {
	DDURL  string `json:"dd_url"`
	Status int    `json:"status"`
	Icon   string `json:"icon"`
}

// DefaultVipCard 默认VIP卡信息
type DefaultVipCard struct {
	PID   int    `json:"pid"`
	PType int    `json:"ptype"`
	Img   string `json:"img"`
}

// ChannelInfo 获取学习圈频道信息
func (s *Service) ChannelInfo(channelID int) (info *ChannelInfo, err error) {

	body, err := s.reqChannelInfo(channelID)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

// ChannelHomepage 获取学习圈频道首页分类
func (s *Service) ChannelHomepage(channelID int) (cats []ChannelHomepageCategory, err error) {

	body, err := s.reqChannelHomepage(channelID)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &cats); err != nil {
		return
	}
	return
}

// ChannelVipInfo 获取学习圈VIP/权限信息
func (s *Service) ChannelVipInfo(channelID int) (info *ChannelVipInfo, err error) {
	body, err := s.reqChannelVipInfo(channelID)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

// ChannelTopicDetail 获取学习圈主题详情
func (s *Service) ChannelTopicDetail(channelID int, topicID int) (topic *ChannelTopicCategory, err error) {
	homepage, err := s.ChannelHomepage(channelID)
	if err != nil {
		return
	}

	// 在所有分类中查找指定主题
	for _, category := range homepage {
		for _, t := range category.List {
			if t.ID == topicID {
				topic = &t
				return topic, nil
			}
		}
	}

	return nil, nil
}

// ChannelTopicItems 获取学习圈主题下的所有条目
func (s *Service) ChannelTopicItems(channelID int, topicID int) (items []ChannelItem, err error) {
	topic, err := s.ChannelTopicDetail(channelID, topicID)
	if err != nil {
		return
	}
	if topic != nil {
		items = topic.Items
	}
	return items, nil
}

// ChannelItemDetail 获取学习圈条目详情（根据enid获取条目信息）
func (s *Service) ChannelItemDetail(channelID int, enid string) (item *ChannelItem, err error) {
	homepage, err := s.ChannelHomepage(channelID)
	if err != nil {
		return
	}

	// 在所有分类和主题中查找指定enid的条目
	for _, category := range homepage {
		for _, topic := range category.List {
			for _, it := range topic.Items {
				if it.EnID == enid {
					item = &it
					return item, nil
				}
			}
		}
	}

	return nil, nil
}
