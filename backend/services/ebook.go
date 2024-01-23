package services

import "github.com/yann0917/dedao-gui/backend/utils"

// Catalog ebook catalog
type Catalog struct {
	Level     int    `json:"level"`
	Text      string `json:"text"`
	Href      string `json:"href"`
	PlayOrder int    `json:"playOrder"`
}

// Press ebook press info
type Press struct {
	Name  string `json:"name"`
	Brief string `json:"brief"`
}

// EbookDetail ebook detail
type EbookDetail struct {
	ID                  int           `json:"id"`
	Title               string        `json:"title"`
	Style               int           `json:"style"`
	Cover               string        `json:"cover"`
	Count               int           `json:"count"`
	Price               string        `json:"price"`
	Status              int           `json:"status"`
	OperatingTitle      string        `json:"operating_title"`
	OtherShareTitle     string        `json:"other_share_title"`
	OtherShareSummary   string        `json:"other_share_summary"`
	AuthorInfo          string        `json:"author_info"`
	BookAuthor          string        `json:"book_author"`
	PublishTime         string        `json:"publish_time"`
	CatalogList         []Catalog     `json:"catalog_list"`
	BookIntro           string        `json:"book_intro"`
	BSpecialPrice       string        `json:"b_special_price"`
	CurrentPrice        string        `json:"current_price"`
	IsBuy               bool          `json:"is_buy"`
	IsTrial             bool          `json:"is_trial"`
	IsTtsSwitch         bool          `json:"is_tts_switch"`
	LogID               string        `json:"log_id"`
	LogType             string        `json:"log_type"`
	OriginalPrice       string        `json:"original_price"`
	AuthorList          []string      `json:"author_list"`
	CanTrialRead        bool          `json:"can_trial_read"`
	TrialReadProportion string        `json:"trial_read_proportion"`
	WithVideo           bool          `json:"with_video"`
	Enid                string        `json:"enid"`
	BOverseasPurchase   int           `json:"b_overseas_purchase"`
	RankName            string        `json:"rank_name"`
	RankNum             int           `json:"rank_num"`
	IsVipBook           int           `json:"is_vip_book"`
	IsOnBookshelf       bool          `json:"is_on_bookshelf"`
	ProductScore        string        `json:"product_score"`
	ReadTime            int           `json:"read_time"`
	ReadNumber          []interface{} `json:"read_number"`
	Press               Press         `json:"press"`
	DoubanScore         string        `json:"douban_score"`
	ClassifyName        string        `json:"classify_name"`
	ClassifyID          int           `json:"classify_id"`
	AddStudylistDdURL   string        `json:"add_studylist_dd_url"`
}

// EbookBlock ebook block
type EbookBlock struct {
	ChapterID   string `json:"chapterId"`
	SectionID   string `json:"sectionID"`
	EndOffset   int    `json:"endOffset"`
	StartOffset int    `json:"startOffset"`
}

// EbookOrders ebook orders
type EbookOrders struct {
	ChapterID  string `json:"chapterId"`
	PathInEpub string `json:"pathInEpub"`
}

// EbookToc ebook toc
type EbookToc struct {
	Href      string `json:"href"`
	Level     int    `json:"level"`
	PlayOrder int    `json:"playOrder"`
	Offset    int    `json:"offset"`
	Text      string `json:"text"`
}

type EbookInfoPage struct {
	Cid         string `json:"cid"`
	EndOffset   int    `json:"end_offset"`
	PageNum     int    `json:"page_num"`
	StartOffset int    `json:"start_offset"`
}

type EbookPage struct {
	IsEnd bool `json:"is_end"`
	Pages []struct {
		BeginOffset           int64  `json:"begin_offset"`
		EndOffset             int64  `json:"end_offset"`
		IsFirst               bool   `json:"is_first"`
		IsLast                bool   `json:"is_last"`
		Svg                   string `json:"svg"`
		ViewHeighToChapterTop int64  `json:"view_heigh_to_chapter_top"`
	} `json:"pages"`
}

// EbookInfo ebook info
type EbookInfo struct {
	BookInfo struct {
		EbookBlock [][]EbookBlock   `json:"block"`
		Orders     []EbookOrders    `json:"orders"`
		Toc        []utils.EbookToc `json:"toc"`
		Pages      []EbookInfoPage  `json:"pages"`
	} `json:"bookInfo"`
}

// EbookVIPInfo ebook vip info
type EbookVIPInfo struct {
	UID                int           `json:"uid"`
	Nickname           string        `json:"nickname"`
	Slogan             string        `json:"slogan"`
	Avatar             string        `json:"avatar"`
	AvatarS            string        `json:"avatar_s"`
	MonthCount         int           `json:"month_count"`
	TotalCount         int           `json:"total_count"`
	WeekCount          int           `json:"week_count"`
	FinishedCount      int           `json:"finished_count"`
	SavePrice          string        `json:"save_price"`
	IsVip              bool          `json:"is_vip"`
	BeginTime          int           `json:"begin_time"`
	EndTime            int           `json:"end_time"`
	EnterpriseEndTime  int           `json:"enterprise_end_time"`
	ExpireTime         int           `json:"expire_time"`
	SurplusTime        int           `json:"surplus_time"`
	IsExpire           bool          `json:"is_expire"`
	CardID             int           `json:"card_id"`
	CardType           int           `json:"card_type"`
	PriceDesc          string        `json:"price_desc"`
	IsBuyMonthDiscount bool          `json:"is_buy_month_discount"`
	MonthDiscountPrice int           `json:"month_discount_price"`
	DdURL              string        `json:"dd_url"`
	ErrTips            string        `json:"err_tips"`
	VStateValue        int           `json:"v_state_value"`
	UpgradeTips        []interface{} `json:"upgrade_tips"`
}

type EbookCommentList struct {
	EbookScore EbookScore `json:"ebook_score"`
	SelfInfo   SelfInfo   `json:"self_info"`
	List       []Comment  `json:"list"`
	Total      int        `json:"total"`
	Book       Book       `json:"book"`
}

type SelfInfo struct {
	Score      int `json:"score"`
	AuditState int `json:"audit_state"`
}

type Book struct {
	BookType int `json:"book_type"`
}

type EbookScore struct {
	Id           int      `json:"id"`
	Pid          int      `json:"pid"`
	Ptype        int      `json:"ptype"`
	Isbn         string   `json:"isbn"`
	AverageScore string   `json:"average_score"`
	ScoreInfo    []string `json:"score_info"`
	CreateTime   string   `json:"create_time"`
	UpdateTime   string   `json:"update_time"`
	Total        string   `json:"total"`
	Status       int      `json:"status"`
	BookStatus   int      `json:"book_status"`
}

type EbookShelfAddResp struct {
	Data struct {
		N     int `json:"n"`
		Count int `json:"count"`
	} `json:"data"`
}

// EbookDetail get ebook detail
func (s *Service) EbookDetail(enid string) (detail *EbookDetail, err error) {

	body, err := s.reqEbookDetail(enid)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

// EbookReadToken get ebook read token
func (s *Service) EbookReadToken(enid string) (t *Token, err error) {
	body, err := s.reqEbookReadToken(enid)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &t); err != nil {
		return
	}
	return
}

// EbookInfo get ebook info
// include book block, book TOC, epubPath etc
func (s *Service) EbookInfo(token string) (info *EbookInfo, err error) {
	body, err := s.reqEbookInfo(token)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

func (s *Service) EbookPages(chapterID, token string, index, count, offset int) (pages *EbookPage, err error) {
	body, err := s.reqEbookPages(chapterID, token, index, count, offset)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &pages); err != nil {
		return
	}
	return
}

// EbookCommentList get ebook comment list
// sort like_count
func (s *Service) EbookCommentList(id, sort string, page, limit int) (list *EbookCommentList, err error) {

	body, err := s.reqEbookCommentList(id, sort, page, limit)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

// EbookShelfAdd add ebook shelf
func (s *Service) EbookShelfAdd(ids []string) (resp *EbookShelfAddResp, err error) {

	body, err := s.reqEbookShelfAdd(ids)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}

// EbookShelfRemove remove from ebook shelf
func (s *Service) EbookShelfRemove(ids []string) (resp *EbookShelfAddResp, err error) {

	body, err := s.reqEbookShelfRemove(ids)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}
