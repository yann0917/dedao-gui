package downloadmgr

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	dbPath, err := defaultDBPath()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&DownloadTask{}, &DownloadTaskLog{}); err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) CreateTask(req CreateTaskRequest) (*DownloadTask, error) {
	if req.MaxRetries <= 0 {
		req.MaxRetries = 3
	}
	var existed DownloadTask
	activeStatuses := []string{StatusPending, StatusQueued, StatusRunning, StatusPaused}
	err := r.db.Where("biz_type = ? AND biz_id = ? AND download_type = ? AND save_dir = ? AND status IN ?",
		req.BizType, req.BizID, req.DownloadType, req.SaveDir, activeStatuses).
		Order("created_at DESC").
		First(&existed).Error
	if err == nil {
		return &existed, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	task := &DownloadTask{
		ID:           uuid.NewString(),
		BizType:      req.BizType,
		BizID:        req.BizID,
		Title:        req.Title,
		EnID:         req.EnID,
		ArticleID:    req.ArticleID,
		DownloadType: req.DownloadType,
		Status:       StatusPending,
		Priority:     req.Priority,
		SaveDir:      req.SaveDir,
		PayloadJSON:  req.PayloadJSON,
		MaxRetries:   req.MaxRetries,
		CurrentName:  req.Title,
	}
	if err = r.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) ListTasks(query ListTaskQuery) (result ListTaskResult, err error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 200 {
		query.PageSize = 20
	}
	tx := r.db.Model(&DownloadTask{})
	if len(query.Status) > 0 {
		tx = tx.Where("status IN ?", query.Status)
	}
	if len(query.BizType) > 0 {
		tx = tx.Where("biz_type IN ?", query.BizType)
	}
	if err = tx.Count(&result.Total).Error; err != nil {
		return
	}
	if err = tx.Order("priority DESC").Order("created_at ASC").
		Offset((query.Page - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&result.List).Error; err != nil {
		return
	}
	result.Page = query.Page
	result.PageSize = query.PageSize
	return
}

func (r *Repository) GetTask(id string) (*DownloadTask, error) {
	var task DownloadTask
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) MarkStartupRecover() error {
	return r.db.Model(&DownloadTask{}).
		Where("status IN ?", []string{StatusQueued, StatusRunning}).
		Updates(map[string]interface{}{
			"status":        StatusPending,
			"error_code":    "RECOVERED",
			"error_message": "任务在应用重启后已恢复为待执行",
			"updated_at":    time.Now(),
		}).Error
}

func (r *Repository) ClaimPendingTasks(limit int) ([]DownloadTask, error) {
	if limit <= 0 {
		limit = 1
	}
	var selected []DownloadTask
	now := time.Now()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("status = ?", StatusPending).
			Where("(next_retry_at IS NULL OR next_retry_at <= ?)", now).
			Order("priority DESC").
			Order("created_at ASC").
			Limit(limit).
			Find(&selected).Error
		if err != nil {
			return err
		}
		for i := range selected {
			res := tx.Model(&DownloadTask{}).
				Where("id = ? AND status = ?", selected[i].ID, StatusPending).
				Updates(map[string]interface{}{
					"status":     StatusQueued,
					"updated_at": time.Now(),
				})
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var claimed []DownloadTask
	for _, task := range selected {
		var latest DownloadTask
		err = r.db.First(&latest, "id = ?", task.ID).Error
		if err != nil {
			return nil, err
		}
		if latest.Status == StatusQueued {
			claimed = append(claimed, latest)
		}
	}
	return claimed, nil
}

func (r *Repository) UpdateTaskStatus(id string, status string, fields map[string]interface{}) error {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["status"] = status
	fields["updated_at"] = time.Now()
	return r.db.Model(&DownloadTask{}).Where("id = ?", id).Updates(fields).Error
}

func (r *Repository) UpdateTaskProgress(id string, progress int, current int, total int, name string) error {
	return r.db.Model(&DownloadTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"progress":     progress,
		"current":      current,
		"total":        total,
		"current_name": name,
		"updated_at":   time.Now(),
	}).Error
}

func (r *Repository) CancelTask(id string) error {
	task, err := r.GetTask(id)
	if err != nil {
		return err
	}
	if IsTerminalStatus(task.Status) {
		return nil
	}
	now := time.Now()
	return r.UpdateTaskStatus(id, StatusCanceled, map[string]interface{}{
		"finished_at":   &now,
		"error_code":    "",
		"error_message": "",
	})
}

func (r *Repository) PauseTask(id string) error {
	task, err := r.GetTask(id)
	if err != nil {
		return err
	}
	if IsTerminalStatus(task.Status) {
		return errors.New("任务已结束，无法暂停")
	}
	return r.UpdateTaskStatus(id, StatusPaused, map[string]interface{}{
		"error_code":    "",
		"error_message": "",
	})
}

func (r *Repository) ResumeTask(id string) error {
	task, err := r.GetTask(id)
	if err != nil {
		return err
	}
	if task.Status != StatusPaused {
		return errors.New("仅暂停中的任务可继续")
	}
	return r.db.Model(&DownloadTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        StatusPending,
		"next_retry_at": nil,
		"error_code":    "",
		"error_message": "",
		"updated_at":    time.Now(),
	}).Error
}

func (r *Repository) RetryTask(id string) error {
	task, err := r.GetTask(id)
	if err != nil {
		return err
	}
	if task.Status != StatusFailed && task.Status != StatusCanceled {
		return errors.New("当前状态不允许重试")
	}
	return r.db.Model(&DownloadTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        StatusPending,
		"retry_count":   0,
		"next_retry_at": nil,
		"error_code":    "",
		"error_message": "",
		"started_at":    nil,
		"finished_at":   nil,
		"updated_at":    time.Now(),
	}).Error
}

func (r *Repository) CreateLog(taskID, level, msg string) error {
	return r.db.Create(&DownloadTaskLog{
		TaskID:  taskID,
		Level:   level,
		Message: msg,
	}).Error
}

func (r *Repository) ClearTasks(clearAll bool) error {
	if clearAll {
		return r.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("1 = 1").Delete(&DownloadTaskLog{}).Error; err != nil {
				return err
			}
			return tx.Where("1 = 1").Delete(&DownloadTask{}).Error
		})
	}

	terminalStatuses := []string{StatusSuccess, StatusFailed, StatusCanceled}
	var taskIDs []string
	if err := r.db.Model(&DownloadTask{}).Where("status IN ?", terminalStatuses).Pluck("id", &taskIDs).Error; err != nil {
		return err
	}
	if len(taskIDs) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_id IN ?", taskIDs).Delete(&DownloadTaskLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", taskIDs).Delete(&DownloadTask{}).Error
	})
}

func defaultDBPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(wd, ".cache", "download")
	if err = os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建下载数据库目录失败: %w", err)
	}
	return filepath.Join(dir, "tasks.db"), nil
}
