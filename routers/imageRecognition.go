package routers

import "github.com/gin-gonic/gin"

func ImageRecognitionRouter(rg *gin.RouterGroup) {
	// 图片识别路由
	rg.POST("/recognize", func(c *gin.Context) {
		// 处理图片识别请求的逻辑
		c.JSON(200, gin.H{
			"message": "Image recognition endpoint",
		})
	})
}
