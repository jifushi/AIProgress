package routers

import "github.com/gin-gonic/gin"

func aiRouter(rg *gin.RouterGroup) {
	rg.GET("/get-sessions")           // 获取会话列表
	rg.POST("/send-message")          // 发送信息
	rg.POST("/session-detail")        // 会话详情,用于获取指定会话的消息记录
	rg.POST("/create-session")        // 创建会话
	rg.POST("/create-stream-session") // 创建流式会话
}
