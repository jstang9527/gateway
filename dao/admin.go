package dao

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jstang9527/gateway/dto"

	"github.com/jstang9527/gateway/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// Admin ...
type Admin struct {
	ID       int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt     string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password string    `json:"password" gorm:"column:password" description:"管理员密码"`
	UpdateAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete int       `json:"is_delete" gorm:"column:is_delete" description:"账户状态"`
}

// TableName 表名
func (a *Admin) TableName() string {
	return "gateway_admin"
}

// Find 数据库表查询
func (a *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginInputParamsCheck 对密码加盐校验
func (a *Admin) LoginInputParamsCheck(c *gin.Context, tx *gorm.DB, params *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := a.Find(c, tx, &Admin{UserName: params.UserName, IsDelete: 0}) //查询用户名为xxx且未被注销的账户
	if err != nil {
		return nil, errors.New("账户不存在或被注销")
	}
	saltPassword := public.GetSaltPassword(adminInfo.Salt, params.Password)
	if saltPassword != adminInfo.Password {
		return nil, errors.New("密码不正确")
	}
	return adminInfo, nil
}
