package services

import "fmt"

// RecentRequest 最近学习接口请求参数
type RecentRequest struct {
	FilterProductType bool   `json:"filter_product_type"`
	MaxID             int64  `json:"max_id"`
	PageSize          int    `json:"page_size"`
	ProductType       string `json:"product_type"`
	UIDHazy           string `json:"uid_hazy"`
}

// RecentResponse 最近学习接口响应
type RecentResponse struct {
	List        []RecentItem `json:"list"`
	Timestamp   int64        `json:"timestamp"`
	CurrentTime int          `json:"current_time"`
	HasMore     bool         `json:"has_more"`
}

// RecentItem 最近学习列表项（只保留界面用到的字段）
type RecentItem struct {
	ProductID     int            `json:"product_id"`
	ProductType   int            `json:"product_type"`
	ProductIDHazy string         `json:"product_id_hazy"`
	TypeName      string         `json:"type_name"`
	Title         string         `json:"title"`
	Author        string         `json:"author"`
	IndexImg      string         `json:"index_img"`
	SquareImg     string         `json:"square_img"`
	Timestamp     int64          `json:"timestamp"`
	LastInfo      string         `json:"last_info"`
	ProgressIntro RecentProgress `json:"progress_intro"`
}

// RecentProgress 最近学习进度
type RecentProgress struct {
	Intro       string `json:"intro"`
	Progress    int    `json:"progress"`
	MaxProgress int    `json:"max_progress"`
	IsFinish    int    `json:"is_finish"`
	Uint        string `json:"uint"`
}

// RecentNameCount 按内容类型聚合的条数
type RecentNameCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// RecentDayCount 每日条数，Date 为 MM-DD
type RecentDayCount struct {
	Date  string            `json:"date"`
	Count int               `json:"count"`
	Types []RecentNameCount `json:"types"` // 按类型细分，顺序与顶层 Types 一致
}

// RecentHourCount 单小时条数，仅非 0 小时
type RecentHourCount struct {
	Hour  int               `json:"hour"`
	Count int               `json:"count"`
	Types []RecentNameCount `json:"types"`
}

// RecentStats 近 N 天学习统计（服务端翻页拉满窗口后聚合）
type RecentStats struct {
	Total      int               `json:"total"`
	WindowDays int               `json:"window_days"`
	ActiveDays int               `json:"active_days"`
	FinishRate int               `json:"finish_rate"` // 百分比整数
	TopType    string            `json:"top_type"`
	Truncated  bool              `json:"truncated"`
	Types      []RecentNameCount `json:"types"`
	Daily      []RecentDayCount  `json:"daily"` // 近 14 天从旧到新，含 0
	Hours      []RecentHourCount `json:"hours"` // 仅非 0 小时
}

// Recent 获取最近学习情况
func (s *Service) Recent(param RecentRequest) (resp *RecentResponse, err error) {
	if param.UIDHazy == "" {
		err = fmt.Errorf("uid_hazy is required")
		return
	}
	if param.PageSize <= 0 {
		param.PageSize = 20
	}

	body, err := s.reqRecent(param)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}
