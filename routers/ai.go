package routers

import (
	"AiProgress/controller/aichat"

	"github.com/gin-gonic/gin"
)

func aiRouter(rg *gin.RouterGroup) {
	rg.GET("/get-sessions", aichat.GetSessions)               // 获取会话列表
	rg.POST("/send-message", aichat.SendMessage)              // 发送信息
	rg.POST("/send-stream-message", aichat.SendStreamMessage) // 发送流式信息
	rg.POST("/session-detail", aichat.SessionDetail)          // 会话详情,用于获取指定会话的消息记录
	rg.POST("/create-session", aichat.CreateSession)          // 创建会话
	// rg.POST("/create-and-send-stream-session", aichat.CreateAndSendStreamSession) // 创建流式会话
	// rg.POST("/create-and-send-session", aichat.CreateAndSendSession)              // 创建并发送会话
	// rg.POST("/delete-session", aichat.DeleteSession)                              // 删除会话
}
