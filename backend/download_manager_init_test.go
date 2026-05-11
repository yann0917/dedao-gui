package backend

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yann0917/dedao-gui/backend/downloadmgr"
)

func TestEnsureDownloadManagerInitializesOnceForConcurrentCalls(t *testing.T) {
	app := NewApp()
	app.Ctx = context.Background()

	var initCalls int32
	app.downloadInitFn = func(ctx context.Context) (*downloadmgr.Repository, *downloadmgr.Manager, error) {
		atomic.AddInt32(&initCalls, 1)
		time.Sleep(20 * time.Millisecond)
		return &downloadmgr.Repository{}, &downloadmgr.Manager{}, nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- app.ensureDownloadManager()
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("ensureDownloadManager returned error: %v", err)
		}
	}
	if got := atomic.LoadInt32(&initCalls); got != 1 {
		t.Fatalf("expected exactly one init call, got %d", got)
	}
	if app.DownloadRepo == nil || app.DownloadManager == nil {
		t.Fatalf("expected download manager resources to be initialized")
	}
	if app.downloadState != downloadManagerStateReady {
		t.Fatalf("expected ready state, got %s", app.downloadState)
	}
}

func TestEnsureDownloadManagerRetriesAfterFailure(t *testing.T) {
	app := NewApp()
	app.Ctx = context.Background()

	var initCalls int32
	app.downloadInitFn = func(ctx context.Context) (*downloadmgr.Repository, *downloadmgr.Manager, error) {
		call := atomic.AddInt32(&initCalls, 1)
		if call == 1 {
			return nil, nil, errors.New("sqlite locked")
		}
		return &downloadmgr.Repository{}, &downloadmgr.Manager{}, nil
	}

	err := app.ensureDownloadManager()
	if err == nil || !strings.Contains(err.Error(), "sqlite locked") {
		t.Fatalf("expected first init failure to be returned, got %v", err)
	}
	if app.downloadState != downloadManagerStateFailed {
		t.Fatalf("expected failed state after first init, got %s", app.downloadState)
	}

	err = app.ensureDownloadManager()
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	if got := atomic.LoadInt32(&initCalls); got != 2 {
		t.Fatalf("expected two init attempts, got %d", got)
	}
	if app.downloadState != downloadManagerStateReady {
		t.Fatalf("expected ready state after retry, got %s", app.downloadState)
	}
}

func TestCreateDownloadTaskReturnsInitError(t *testing.T) {
	app := NewApp()
	app.Ctx = context.Background()
	app.downloadInitFn = func(ctx context.Context) (*downloadmgr.Repository, *downloadmgr.Manager, error) {
		return nil, nil, errors.New("创建下载数据库目录失败")
	}

	_, err := app.CreateDownloadTask(downloadmgr.CreateTaskRequest{
		BizType:      downloadmgr.BizTypeCourse,
		BizID:        1,
		DownloadType: 1,
	})
	if err == nil {
		t.Fatal("expected init error, got nil")
	}
	if !strings.Contains(err.Error(), "创建下载数据库目录失败") {
		t.Fatalf("expected real init error to bubble up, got %v", err)
	}
}
