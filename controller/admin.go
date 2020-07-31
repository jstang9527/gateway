package controller

import (
	"encoding/json"
	"fmt"

	"github.com/e421083458/golang_common/lib"

	"github.com/jstang9527/gateway/dao"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
	"github.com/jstang9527/gateway/public"
)

// AdminController ...
type AdminController struct{}

// AdminRegister 登录控制器
func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admin_info", admin.Admin)
	group.POST("/change_pwd", admin.ChangePwd)
}

// Admin godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (a *AdminController) Admin(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey) //获得结构体的interface
	sessInfoStr := fmt.Sprint(sessInfo)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfoStr), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 3000, err)
		return
	}
	// 1. 读取sessionKey对应的json转换成结构体
	// 2. 取数据进行封装结构体
	// 4. 返回信息
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		UserName:     adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://zan71.com/cdn-img/icon/avatar/admin.icon",
		Introduction: "super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept json
// @Produce json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *AdminController) ChangePwd(c *gin.Context) {
	//1. 请求参数(密码)初步校验(必填)
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3001, err)
		return
	}
	//2. session读取用户信息结构体 adminSessionInfo
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey) //获得结构体的interface
	sessInfoStr := fmt.Sprint(sessInfo)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfoStr), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 3002, err)
		return
	}
	//3. sessInfo.ID读取数据库信息 adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 3003, err)
		return
	}
	adminInfo := &dao.Admin{}
	fmt.Println(adminSessionInfo.UserName)
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 3004, err)
		return
	}
	//4. params.password+admininfo.salt sha256 => saltpassword
	saltPassword := public.GetSaltPassword(adminInfo.Salt, params.Password)
	//5. saltpassword ==> adminInfo.password 执行数据库保存
	adminInfo.Password = saltPassword
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 3005, err)
		return
	}
	//4. 返回信息
	out := fmt.Sprintf("update password success in %v", adminInfo.UpdatedAt)
	middleware.ResponseSuccess(c, out)
}
