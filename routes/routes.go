package routes

import (
	"web_app/controller"
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	//设置不同模式控制不同的日志输出
	if mode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")

	//1. 注册处理
	v1.POST("/signup", controller.SignUpHandler)

	//2. 登录处理
	v1.POST("/login", controller.LoginHandler)

	//3. 使用JWT认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		//社区分类和细节查询
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		//根据id访问帖子
		v1.GET("/:id", controller.GetPostDetailHandler)
		//查看所有帖子
		v1.GET("/posts/", controller.GetPostListHandler)
		//查看所有帖子（根据分页、排序需求）
		v1.GET("/posts2/", controller.GetPostListHandler2)

		//对帖子投票（点赞）
		v1.POST("/vote", controller.PostVoteHandler)
		//创建帖子请求
		v1.POST("/post", controller.CreatePostHandler)
	}

	return r
}
