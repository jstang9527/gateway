package router

import (
	"log"

	"github.com/jstang9527/gateway/controller"

	"github.com/jstang9527/gateway/middleware"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter ...
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//登录
	adminLoginRouter := router.Group("/admin_login")
	store, err := sessions.NewRedisStore(10, "tcp", "127.0.0.1:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("sessions.NewRedisStore err: %v", err)
	}
	adminLoginRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}
	return router
}