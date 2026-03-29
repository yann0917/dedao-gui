package downloadmgr

import "time"

const (
	BizTypeCourse = "course"
	BizTypeOdob   = "odob"
	BizTypeEbook  = "ebook"
)

const (
	StatusPending  = "pending"
	StatusQueued   = "queued"
	StatusRunning  = "running"
	StatusPaused   = "paused"
	StatusSuccess  = "success"
	StatusFailed   = "failed"
	StatusCanceled = "canceled"
)

type DownloadTask struct {
	ID           string     `gorm:"primaryKey;type:varchar(64)" json:"id"`
	BizType      string     `gorm:"type:varchar(32);index:idx_tasks_status_priority_created,priority:1;index" json:"bizType"`
	BizID        int        `gorm:"index" json:"bizId"`
	Title        string     `gorm:"type:text" json:"title"`
	EnID         string     `gorm:"type:varchar(255)" json:"enId"`
	ArticleID    int        `json:"articleId"`
	DownloadType int        `json:"downloadType"`
	Status       string     `gorm:"type:varchar(32);index:idx_tasks_status_priority_created,priority:2;index" json:"status"`
	Priority     int        `gorm:"default:0;index:idx_tasks_status_priority_created,priority:3" json:"priority"`
	Progress     int        `gorm:"default:0" json:"progress"`
	Current      int        `gorm:"default:0" json:"current"`
	Total        int        `gorm:"default:0" json:"total"`
	CurrentName  string     `gorm:"type:text" json:"currentName"`
	SaveDir      string     `gorm:"type:text" json:"saveDir"`
	PayloadJSON  string     `gorm:"type:text" json:"payloadJson"`
	RetryCount   int        `gorm:"default:0" json:"retryCount"`
	MaxRetries   int        `gorm:"default:3" json:"maxRetries"`
	NextRetryAt  *time.Time `gorm:"index" json:"nextRetryAt"`
	ErrorCode    string     `gorm:"type:varchar(64)" json:"errorCode"`
	ErrorMessage string     `gorm:"type:text" json:"errorMessage"`
	StartedAt    *time.Time `json:"startedAt"`
	FinishedAt   *time.Time `json:"finishedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type DownloadTaskLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    string    `gorm:"type:varchar(64);index" json:"taskId"`
	Level     string    `gorm:"type:varchar(16);index" json:"level"`
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateTaskRequest struct {
	BizType      string `json:"bizType"`
	BizID        int    `json:"bizId"`
	Title        string `json:"title"`
	EnID         string `json:"enId"`
	ArticleID    int    `json:"articleId"`
	DownloadType int    `json:"downloadType"`
	SaveDir      string `json:"saveDir"`
	Priority     int    `json:"priority"`
	PayloadJSON  string `json:"payloadJson"`
	MaxRetries   int    `json:"maxRetries"`
}

type ListTaskQuery struct {
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Status   []string `json:"status"`
	BizType  []string `json:"bizType"`
}

type ListTaskResult struct {
	List     []DownloadTask `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

func IsTerminalStatus(status string) bool {
	return status == StatusSuccess || status == StatusFailed || status == StatusCanceled
}
