package dto

// PanelDataOutput 管理员信息输出
type PanelDataOutput struct {
	ServiceNum  int64 `json:"service_num"`       //服务数量
	AppNum      int64 `json:"app_num"`           //租户数
	CurrentQPS  int64 `json:"current_qps"`       //当前QPS
	TodayReqNum int64 `json:"today_request_num"` //今日请求量
}

// PanelSrvStatItemOutput ...
type PanelSrvStatItemOutput struct {
	Name     string `json:"name"`
	LoadType int    `json:"load_type"`
	Value    int64  `json:"value"`
}

// PanelSrvStatOutput 管理员信息输出
type PanelSrvStatOutput struct {
	Legend []string                 `json:"legend"`
	Data   []PanelSrvStatItemOutput `json:"data"`
}
