package app

import (
	"errors"
	"sort"
	"time"

	"github.com/yann0917/dedao-gui/backend/config"
	"github.com/yann0917/dedao-gui/backend/services"
)

// RecentList 最近学习记录。
// 翻页为游标式：下一页 maxID 传上一页最后一条的 timestamp（毫秒），实测两页无重叠。
func RecentList(maxID int64, pageSize int) (list *services.RecentResponse, err error) {
	uid := config.Instance.ActiveUser().UIDHazy
	if uid == "" {
		err = errors.New("未登录")
		return
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	ser := getService()
	list, err = ser.Recent(services.RecentRequest{
		FilterProductType: true,
		MaxID:             maxID,
		PageSize:          pageSize,
		UIDHazy:           uid,
	})
	return
}

const (
	recentStatsWindowDays = 30 // 统计窗口：近 30 天
	recentStatsMaxRecords = 500
	recentStatsDailyDays  = 14 // 每日条数图表的天数
)

// RecentStats 近 30 天学习统计：服务端翻页拉满窗口后聚合，
// 记录按学习时间倒序，遇到早于窗口起点的整页即可停止。
func RecentStats() (stats *services.RecentStats, err error) {
	uid := config.Instance.ActiveUser().UIDHazy
	if uid == "" {
		err = errors.New("未登录")
		return
	}
	ser := getService()

	now := time.Now()
	// 自然日语义：含今天共 30 天，避免滚动窗口横跨 31 个日历日
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	cutoff := todayStart.AddDate(0, 0, -(recentStatsWindowDays - 1))
	byType := map[string]int{}
	dayTypes := map[string]map[string]int{} // MM-DD -> 类型 -> 条数
	hourTypes := map[int]map[string]int{}
	activeDays := map[string]bool{}
	finished := 0
	collected := 0

	maxID := int64(0)
	truncated := false
	for {
		resp, err := ser.Recent(services.RecentRequest{
			FilterProductType: true,
			MaxID:             maxID,
			PageSize:          100,
			UIDHazy:           uid,
		})
		if err != nil {
			return nil, err
		}
		if len(resp.List) == 0 {
			break
		}

		oldest := time.Now()
		for _, it := range resp.List {
			if it.Timestamp <= 0 {
				continue
			}
			t := time.UnixMilli(it.Timestamp)
			if t.Before(oldest) {
				oldest = t
			}
			if t.Before(cutoff) {
				continue
			}
			collected++
			if it.ProgressIntro.IsFinish == 1 {
				finished++
			}
			typeName := it.TypeName
			if typeName == "" {
				typeName = "其他"
			}
			byType[typeName]++
			day := t.Format("01-02")
			if dayTypes[day] == nil {
				dayTypes[day] = map[string]int{}
			}
			dayTypes[day][typeName]++
			if hourTypes[t.Hour()] == nil {
				hourTypes[t.Hour()] = map[string]int{}
			}
			hourTypes[t.Hour()][typeName]++
			activeDays[day] = true
		}

		if collected >= recentStatsMaxRecords {
			// 记录数达到上限：仅当流未结束且本页最旧记录仍在窗口内时，才视为有截断
			truncated = resp.HasMore && !oldest.Before(cutoff)
			break
		}
		if !resp.HasMore || oldest.Before(cutoff) {
			break
		}
		maxID = resp.List[len(resp.List)-1].Timestamp
		if maxID <= 0 {
			break
		}
	}

	types := make([]services.RecentNameCount, 0, len(byType))
	topType := ""
	for name, count := range byType {
		types = append(types, services.RecentNameCount{Name: name, Count: count})
		if count > byType[topType] {
			topType = name
		}
	}
	sort.Slice(types, func(i, j int) bool {
		if types[i].Count != types[j].Count {
			return types[i].Count > types[j].Count
		}
		return types[i].Name < types[j].Name
	})
	typeOrder := make(map[string]int, len(types))
	for i, t := range types {
		typeOrder[t.Name] = i
	}

	// orderedTypes 按顶层 types 的顺序输出细分，前端堆叠时各图分段顺序一致
	orderedTypes := func(m map[string]int) []services.RecentNameCount {
		out := make([]services.RecentNameCount, 0, len(m))
		for name, count := range m {
			out = append(out, services.RecentNameCount{Name: name, Count: count})
		}
		sort.Slice(out, func(i, j int) bool { return typeOrder[out[i].Name] < typeOrder[out[j].Name] })
		return out
	}

	daily := make([]services.RecentDayCount, 0, recentStatsDailyDays)
	for i := recentStatsDailyDays - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i).Format("01-02")
		m := dayTypes[day]
		total := 0
		for _, count := range m {
			total += count
		}
		daily = append(daily, services.RecentDayCount{Date: day, Count: total, Types: orderedTypes(m)})
	}

	hours := make([]services.RecentHourCount, 0, 24)
	for h := 0; h < 24; h++ {
		m := hourTypes[h]
		if len(m) == 0 {
			continue
		}
		total := 0
		for _, count := range m {
			total += count
		}
		hours = append(hours, services.RecentHourCount{Hour: h, Count: total, Types: orderedTypes(m)})
	}

	finishRate := 0
	if collected > 0 {
		finishRate = finished * 100 / collected
	}

	stats = &services.RecentStats{
		Total:      collected,
		WindowDays: recentStatsWindowDays,
		ActiveDays: len(activeDays),
		FinishRate: finishRate,
		TopType:    topType,
		Truncated:  truncated,
		Types:      types,
		Daily:      daily,
		Hours:      hours,
	}
	return
}
