package routers

import (
	"AiProgress/middeware/jwt"

	"github.com/gin-gonic/gin"
)

func RouteController() *gin.Engine {
	r := gin.Default()
	enterRouter := r.Group("/api")
	{
		UserGroup := enterRouter.Group("/user")
		userRouter(UserGroup)
	}
	{
		AIChatGroup := enterRouter.Group("/AIChat")
		AIChatGroup.Use(jwt.Verify()) // 使用前需要jwt认证
		aiRouter(AIChatGroup)
	}
	{
		// 图片识别路由
		ImageGroup := enterRouter.Group("/ImageRecognition")
		ImageGroup.Use(jwt.Verify()) // 使用前需要jwt认证
		ImageRecognitionRouter(ImageGroup)
	}
	{
		// comfyui的文生图路由，或者图片的二次渲染
	}
	return r
}
