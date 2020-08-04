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
	// 1.从session中找admin的session信息并转成session结构体
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	sessInfoStr := fmt.Sprint(sessInfo)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfoStr), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 3000, err)
		return
	}
	// 2.从session结构体中取admin基本信息封装输出结构体
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		UserName:     adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://zan71.com/cdn-img/icon/avatar/tx.gif",
		Introduction: "super administrator",
		Roles:        []string{"admin"},
	}
	// + 扩展: 头像、角色、描述等信息可以查数据库再进行封装，以上为写死
	// 3.返回输出结构体
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
	inputParams := &dto.ChangePwdInput{}
	if err := inputParams.BindValidParam(c); err != nil {
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
	//3. sessInfo.ID读取数据库用户信息 adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 3003, err)
		return
	}
	adminInfo := &dao.Admin{}
	err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 3004, err)
		return
	}
	//4. params.password+admininfo.salt sha256 => saltpassword
	saltPassword := public.GetSaltPassword(adminInfo.Salt, inputParams.Password)
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
