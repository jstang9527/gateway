package controller

import (
	"encoding/json"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dao"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
	"github.com/jstang9527/gateway/public"
)

// AdminLoginController  登录控制器结果体
type AdminLoginController struct{}

// AdminLoginRegister 登录控制器
func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (a *AdminLoginController) AdminLogin(c *gin.Context) {
	//1. 请求参数用户名、密码初步校验(必填)
	inputParams := &dto.AdminLoginInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//2. 账户密码正确性校验
	//2.1 依据inputParams.UserName从DB取管理员信息admininfo
	//2.2 admininfo.salt+parmas.Password sha256 =>saltpassword
	//2.3 判断saltpassword 是否等于admininfo的password
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &dao.Admin{}
	err = admin.LoginInputParamsCheck(c, tx, inputParams)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 3. 保持会话
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.ID,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	// 4. 返回信息
	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}

// AdminLoginOut godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (a *AdminLoginController) AdminLoginOut(c *gin.Context) {
	// 1. 对gin的session进行对于key的删除
	// 1.1 拿到session对象
	sess := sessions.Default(c)
	// 1.2 从对象中删除对于key
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()

	// 2. 返回信息
	out := "logout success."
	middleware.ResponseSuccess(c, out)
}
