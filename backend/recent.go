package backend

import (
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
)

// RecentList 最近学习记录，游标翻页：下一页 maxID 传上一页最后一条的 timestamp
func (a *App) RecentList(maxID int64, pageSize int) (list *services.RecentResponse, err error) {
	list, err = app.RecentList(maxID, pageSize)
	return
}

// RecentStats 近 30 天学习统计（服务端聚合）
func (a *App) RecentStats() (stats *services.RecentStats, err error) {
	stats, err = app.RecentStats()
	return
}
