package dto

import (
	"time"
)

// TestItem ...
type TestItem struct {
	ProjectID   string        //隶属id
	ProjectName string        //项目名称
	ProjectAddr string        //项目地址
	SnapshotIP  string        //部署机器ip
	FuncName    string        //功能名
	Message     string        //测试结果
	Screenshot  string        //测试截图url
	Priority    string        //功能项重要性
	Status      string        //告警程度
	Duration    time.Duration //测试耗时
}

func statusBecomeQuota(level int) string {
	switch level {
	case 1:
		return "Pass"
	case 2:
		return "Warning"
	case 3:
		return "Error"
	default:
		return "Unknown"
	}
}
func priorityBecomeQuota(level int) string {
	switch level {
	case 2:
		return "Nomal"
	case 3:
		return "High"
	default:
		return "Low"
	}
}

// NewTestItem ...
func NewTestItem(pid, pname, paddr, webaddr, fname, msg, shot string, prio, status int, dua time.Duration) *TestItem {
	return &TestItem{
		ProjectID:   pid,
		ProjectName: pname,
		ProjectAddr: paddr,
		SnapshotIP:  webaddr,
		FuncName:    fname,
		Message:     msg,
		Screenshot:  shot,
		Priority:    priorityBecomeQuota(prio),
		Status:      statusBecomeQuota(status),
		Duration:    dua,
	}
}
