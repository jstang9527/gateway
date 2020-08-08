package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// ServiceAddHTTPInput ...
type ServiceAddHTTPInput struct {
	ServiceName            string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"`       //服务名
	ServiceDesc            string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`           //服务描述
	RuleType               int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                            //接入类型
	Rule                   string `json:"rule" form:"rule" comment:"接入路径: 域名或者前缀" example:"" validate:"required,valid_rule"`                      //接入路径
	NeedHTTPS              int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                       //是否支持https
	NeedStripURI           int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`           //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"" validate:"max=1,min=0"`         //是否支持websocket
	URLRewrite             string `json:"url_rewrite" form:"url_rewrite" comment:"URL重写" example:"" validate:"valid_url_rewrite"`                 //URL重写
	HeaderTransfor         string `json:"header_transfor" form:"header_transfor" comment:"Header头转换" example:"" validate:"valid_header_transfor"` //Header头转换
	OpenAuth               int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                          //是否开启权限
	BlackList              string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`                                    //黑名单ip
	WhiteList              string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`                                    //白名单ip
	ClientIPFlowLimit      int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`           //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端ip限流" example:"" validate:"min=0"`             //服务端ip限流
	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`                          //轮询方式
	IPList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"required,valid_iplist"`                      //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weightlist"`          //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"connect_timeout" comment:"建立连接超时(秒)" example:"" validate:"min=0"`        //建立连接超时(秒)
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"header_timeout" comment:"获取header超时(秒)" example:"" validate:"min=0"`      //获取header超时(秒)
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"idle_timeout" comment:"连接最大空闲时间(秒)" example:"" validate:"min=0"`            //连接最大空闲时间(秒)
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"max_idle" comment:"最大空闲连接数" example:"" validate:"min=0"`                        //最大空闲连接数
}

// ServiceUpdateHTTPInput ...
type ServiceUpdateHTTPInput struct {
	ID                     int64  `json:"id" form:"id" comment:"服务id" example:"62" validate:"required,min=1"`                                                     //id
	ServiceName            string `json:"service_name" form:"service_name" comment:"服务名" example:"test_http_service_indb" validate:"required,valid_service_name"` //服务名
	ServiceDesc            string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"test_http_service_indb" validate:"required,max=255,min=1"`     //服务描述
	RuleType               int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                                            //接入类型
	Rule                   string `json:"rule" form:"rule" comment:"接入路径: 域名或者前缀" example:"/test_http_service_indb" validate:"required,valid_rule"`               //接入路径
	NeedHTTPS              int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                                       //是否支持https
	NeedStripURI           int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`                           //启用strip_uri
	NeedWebsocket          int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"" validate:"max=1,min=0"`                         //是否支持websocket
	URLRewrite             string `json:"url_rewrite" form:"url_rewrite" comment:"URL重写" example:"" validate:"valid_url_rewrite"`                                 //URL重写
	HeaderTransfor         string `json:"header_transfor" form:"header_transfor" comment:"Header头转换" example:"" validate:"valid_header_transfor"`                 //Header头转换
	OpenAuth               int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                                          //是否开启权限
	BlackList              string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`                                                    //黑名单ip
	WhiteList              string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`                                                    //白名单ip
	ClientIPFlowLimit      int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`                           //客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端ip限流" example:"" validate:"min=0"`                             //服务端ip限流
	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`                                          //轮询方式
	IPList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_iplist"`                          //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"50" validate:"required,valid_weightlist"`                        //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"connect_timeout" comment:"建立连接超时(秒)" example:"" validate:"min=0"`                        //建立连接超时(秒)
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"header_timeout" comment:"获取header超时(秒)" example:"" validate:"min=0"`                      //获取header超时(秒)
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"idle_timeout" comment:"连接最大空闲时间(秒)" example:"" validate:"min=0"`                            //连接最大空闲时间(秒)
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"max_idle" comment:"最大空闲连接数" example:"" validate:"min=0"`                                        //最大空闲连接数
}

// ServiceDeleteInput 服务删除
type ServiceDeleteInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` // 服务ID
}

// ServiceListInput 服务分页查询
type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` // 每页条数
}

// BindValidParam 校验新增参数,绑定结构体,校验参数
func (s *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// BindValidParam 校验新增参数,绑定结构体,校验参数
func (s *ServiceUpdateHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// BindValidParam 校验请求方法,绑定结构体,校验参数
func (s *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// BindValidParam 校验删除参数,绑定结构体,校验参数
func (s *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// ServiceListOutput ...
type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

// ServiceListItemOutput 对象item信息
type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     //id
	ServiceName string `json:"service_name" form:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` //服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       //类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` //服务地址
	QPS         int64  `json:"qps" form:"qps"`                   //qps
	QPD         int64  `json:"qpd" form:"qpd"`                   //qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     //节点数
}

// ServiceStatOutput ...
type ServiceStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" example:"" validate:""`         //今日流量
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" example:"" validate:""` //昨日流量
}

// ServiceAddTCPInput ...
type ServiceAddTCPInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"header头转换" validate:""`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IPList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

// GetValidParams ...
func (params *ServiceAddTCPInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceUpdateTCPInput ...
type ServiceUpdateTCPInput struct {
	ID                int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IPList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

// GetValidParams ...
func (params *ServiceUpdateTCPInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceAddGrpcInput ...
type ServiceAddGrpcInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IPList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

// GetValidParams ...
func (params *ServiceAddGrpcInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceUpdateGrpcInput ...
type ServiceUpdateGrpcInput struct {
	ID                int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IPList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

// GetValidParams ...
func (params *ServiceUpdateGrpcInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
