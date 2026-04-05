package downloadmgr

import (
	"context"
	"errors"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
)

type progressReporter struct {
	taskID string
	repo   *Repository
	ctx    context.Context
}

func (p *progressReporter) Report(progress int, current int, total int, currentName string) {
	_ = p.repo.UpdateTaskProgress(p.taskID, progress, current, total, currentName)
	runtime.EventsEmit(p.ctx, "download:task:update", map[string]interface{}{
		"taskId":      p.taskID,
		"progress":    progress,
		"current":     current,
		"total":       total,
		"currentName": currentName,
	})
}

type Executor interface {
	Execute(ctx context.Context, task *DownloadTask, reporter *progressReporter) error
}

type Manager struct {
	ctx          context.Context
	repo         *Repository
	queue        chan DownloadTask
	workerCount  int
	pollInterval time.Duration
	executors    map[string]Executor
	stopCh       chan struct{}
	stopOnce     sync.Once
	wg           sync.WaitGroup
	runningMu    sync.Mutex
	running      map[string]context.CancelFunc
}

type Config struct {
	WorkerCount  int
	PollInterval time.Duration
}

func NewManager(ctx context.Context, repo *Repository, cfg Config) *Manager {
	workerCount := cfg.WorkerCount
	if workerCount <= 0 {
		workerCount = 1
	}
	pollInterval := cfg.PollInterval
	if pollInterval <= 0 {
		pollInterval = time.Second
	}
	return &Manager{
		ctx:          ctx,
		repo:         repo,
		queue:        make(chan DownloadTask, 100),
		workerCount:  workerCount,
		pollInterval: pollInterval,
		executors:    map[string]Executor{},
		stopCh:       make(chan struct{}),
		running:      map[string]context.CancelFunc{},
	}
}

func (m *Manager) RegisterExecutor(bizType string, executor Executor) {
	m.executors[bizType] = executor
}

func (m *Manager) Start() error {
	if err := m.repo.MarkStartupRecover(); err != nil {
		return err
	}
	m.wg.Add(1)
	go m.dispatchLoop()
	for i := 0; i < m.workerCount; i++ {
		m.wg.Add(1)
		go m.workerLoop()
	}
	return nil
}

func (m *Manager) Stop() {
	m.stopOnce.Do(func() {
		close(m.stopCh)
	})
	m.wg.Wait()
}

func (m *Manager) dispatchLoop() {
	defer m.wg.Done()
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			tasks, err := m.repo.ClaimPendingTasks(m.workerCount * 2)
			if err != nil {
				continue
			}
			for _, task := range tasks {
				select {
				case m.queue <- task:
					m.emitStatus(task.ID, StatusQueued, 0, "", "")
				case <-m.stopCh:
					return
				}
			}
		}
	}
}

func (m *Manager) workerLoop() {
	defer m.wg.Done()
	for {
		select {
		case <-m.stopCh:
			return
		case task := <-m.queue:
			m.execute(task)
		}
	}
}

func (m *Manager) execute(task DownloadTask) {
	executor, ok := m.executors[task.BizType]
	if !ok {
		_ = m.failTask(task, "TASK_EXECUTOR_NOT_FOUND", "未找到对应下载执行器")
		return
	}

	startedAt := time.Now()
	_ = m.repo.UpdateTaskStatus(task.ID, StatusRunning, map[string]interface{}{
		"started_at":    &startedAt,
		"error_code":    "",
		"error_message": "",
	})
	m.emitStatus(task.ID, StatusRunning, 0, "", "")

	ctx, cancel := context.WithCancel(m.ctx)
	m.runningMu.Lock()
	m.running[task.ID] = cancel
	m.runningMu.Unlock()
	defer func() {
		m.runningMu.Lock()
		delete(m.running, task.ID)
		m.runningMu.Unlock()
	}()

	reporter := &progressReporter{
		taskID: task.ID,
		repo:   m.repo,
		ctx:    m.ctx,
	}
	err := executor.Execute(ctx, &task, reporter)
	if err == nil {
		finishedAt := time.Now()
		_ = m.repo.UpdateTaskStatus(task.ID, StatusSuccess, map[string]interface{}{
			"progress":      100,
			"finished_at":   &finishedAt,
			"error_code":    "",
			"error_message": "",
		})
		m.emitStatus(task.ID, StatusSuccess, 100, "", "")
		return
	}
	if errors.Is(err, context.Canceled) {
		latest, e := m.repo.GetTask(task.ID)
		if e == nil && latest.Status == StatusPaused {
			m.emitStatus(task.ID, StatusPaused, task.Progress, "", "")
			return
		}
		finishedAt := time.Now()
		_ = m.repo.UpdateTaskStatus(task.ID, StatusCanceled, map[string]interface{}{
			"finished_at":   &finishedAt,
			"error_code":    "TASK_CANCELED",
			"error_message": "任务已取消",
		})
		m.emitStatus(task.ID, StatusCanceled, task.Progress, "TASK_CANCELED", "任务已取消")
		return
	}
	_ = m.failTask(task, "TASK_INTERNAL_ERROR", err.Error())
}

func (m *Manager) failTask(task DownloadTask, code, msg string) error {
	nextRetry := time.Now().Add(nextBackoff(task.RetryCount))
	retryCount := task.RetryCount + 1
	if retryCount > task.MaxRetries {
		finishedAt := time.Now()
		err := m.repo.UpdateTaskStatus(task.ID, StatusFailed, map[string]interface{}{
			"retry_count":   retryCount,
			"finished_at":   &finishedAt,
			"error_code":    code,
			"error_message": msg,
		})
		m.emitStatus(task.ID, StatusFailed, task.Progress, code, msg)
		return err
	}
	err := m.repo.UpdateTaskStatus(task.ID, StatusPending, map[string]interface{}{
		"retry_count":   retryCount,
		"next_retry_at": &nextRetry,
		"error_code":    code,
		"error_message": msg,
	})
	m.emitStatus(task.ID, StatusPending, task.Progress, code, msg)
	return err
}

func (m *Manager) CancelTask(taskID string) error {
	m.runningMu.Lock()
	cancel, ok := m.running[taskID]
	m.runningMu.Unlock()
	if ok {
		cancel()
		return nil
	}
	return m.repo.CancelTask(taskID)
}

func (m *Manager) PauseTask(taskID string) error {
	m.runningMu.Lock()
	cancel, ok := m.running[taskID]
	m.runningMu.Unlock()
	if ok {
		cancel()
	}
	return m.repo.PauseTask(taskID)
}

func (m *Manager) ResumeTask(taskID string) error {
	return m.repo.ResumeTask(taskID)
}

func (m *Manager) ClearTasks(clearAll bool) error {
	if clearAll {
		m.runningMu.Lock()
		cancels := make([]context.CancelFunc, 0, len(m.running))
		for _, cancel := range m.running {
			cancels = append(cancels, cancel)
		}
		m.runningMu.Unlock()
		for _, cancel := range cancels {
			cancel()
		}
	}
	return m.repo.ClearTasks(clearAll)
}

func (m *Manager) emitStatus(taskID, status string, progress int, errorCode, errorMessage string) {
	runtime.EventsEmit(m.ctx, "download:task:update", map[string]interface{}{
		"taskId":       taskID,
		"status":       status,
		"progress":     progress,
		"errorCode":    errorCode,
		"errorMessage": errorMessage,
		"updatedAt":    time.Now().Format(time.RFC3339),
	})
}

func nextBackoff(retryCount int) time.Duration {
	switch retryCount {
	case 0:
		return time.Minute
	case 1:
		return 2 * time.Minute
	case 2:
		return 5 * time.Minute
	default:
		return 10 * time.Minute
	}
}

type CourseExecutor struct {
	ctx context.Context
}

func NewCourseExecutor(ctx context.Context) *CourseExecutor {
	return &CourseExecutor{ctx: ctx}
}

func (e *CourseExecutor) Execute(ctx context.Context, task *DownloadTask, reporter *progressReporter) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	app.OutputDir = task.SaveDir
	d := app.CourseDownload{
		Ctx:          ctx,
		ID:           task.BizID,
		AID:          task.ArticleID,
		EnId:         task.EnID,
		DownloadType: task.DownloadType,
		ProgressCB: func(total, current, pct int, name string) {
			reporter.Report(pct, current, total, name)
		},
	}
	err := d.Download()
	if err != nil {
		return err
	}
	reporter.Report(100, 1, 1, task.Title)
	return nil
}

type EbookExecutor struct {
	ctx context.Context
}

func NewEbookExecutor(ctx context.Context) *EbookExecutor {
	return &EbookExecutor{ctx: ctx}
}

func (e *EbookExecutor) Execute(ctx context.Context, task *DownloadTask, reporter *progressReporter) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	app.OutputDir = task.SaveDir
	d := app.EBookDownload{
		Ctx:          ctx,
		ID:           task.BizID,
		EnID:         task.EnID,
		DownloadType: task.DownloadType,
		ProgressCB: func(total, current, pct int, name string) {
			reporter.Report(pct, current, total, name)
		},
	}
	err := d.Download()
	if err != nil {
		return err
	}
	reporter.Report(100, 1, 1, task.Title)
	return nil
}

type OdobPayload struct {
	Data services.Course `json:"data"`
}

type OdobExecutor struct {
	ctx context.Context
}

func NewOdobExecutor(ctx context.Context) *OdobExecutor {
	return &OdobExecutor{ctx: ctx}
}

func (e *OdobExecutor) Execute(ctx context.Context, task *DownloadTask, reporter *progressReporter) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	var payload OdobPayload
	if task.PayloadJSON == "" {
		return errors.New("odob 任务缺少 payload")
	}
	if err := jsoniter.UnmarshalFromString(task.PayloadJSON, &payload); err != nil {
		return err
	}
	app.OutputDir = task.SaveDir
	d := app.OdobDownload{
		Ctx:          ctx,
		ID:           task.BizID,
		DownloadType: task.DownloadType,
		Data:         &payload.Data,
		ProgressCB: func(total, current, pct int, name string) {
			reporter.Report(pct, current, total, name)
		},
	}
	err := d.Download()
	if err != nil {
		return err
	}
	reporter.Report(100, 1, 1, task.Title)
	return nil
}
