package dao

// ServiceDetail 多表查询后形成的服务详细信息结构体
type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HTTPRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TCPRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}
