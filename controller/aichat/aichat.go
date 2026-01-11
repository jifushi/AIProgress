package aichat

import (
	"log"
	"net/http"

	common_code "AiProgress/common/code"
	"AiProgress/controller"
	"AiProgress/model"
	service_aichat "AiProgress/service/aichat"

	"github.com/gin-gonic/gin"
)

type (
	GetSessionsRequest  struct{}
	GetSessionsResponse struct {
		Sessions []model.SessionInfo
		controller.Response
	}
	SendMessageRequest struct {
		ModelType string `json:"mode_type"`
		SessionID string `json:"session_id"`
		Content   string `json:"content"`
	}
	SendMessageResponse struct {
		Content string `json:"content"`
		controller.Response
	}
	SendStreamMessageRequest struct {
		ModelType string `json:"mode_type"`
		SessionID string `json:"session_id"`
		Content   string `json:"content"`
	}
	SendStreamMessageResponse struct {
		Content string `json:"content"`
		controller.Response
	}
	SessionDetailRequest struct {
		SessionID string `json:"session_id"`
	}
	SessionDetailResponse struct {
		History []*model.History `json:"history"`
		controller.Response
	}
	CreateSessionRequest  struct{}
	CreateSessionResponse struct {
		SessionID string `json:"session_id"`
		controller.Response
	}
	CreateAndSendStreamSessionRequest  struct{}
	CreateAndSendStreamSessionResponse struct {
		SessionID string `json:"session_id"`
		controller.Response
	}
	CreateAndSendSessionRequest struct {
		Content string `json:"content"`
	}
	CreateAndSendSessionResponse struct {
		Content string `json:"content"`
		controller.Response
	}
	DeleteSessionRequest struct {
		SessionID string `json:"session_id"`
	}
	DeleteSessionResponse struct {
		controller.Response
	}
)

func GetSessions(ctx *gin.Context) {
	var username string
	// var req GetSessionsRequest
	var res GetSessionsResponse
	username = ctx.GetString("username")
	if username == "" {
		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
		return
	}
	sessions, code := service_aichat.GetSessions(username)
	if code != common_code.CodeSuccess {
		ctx.JSON(200, res.CodeOf(code))
		return
	}
	res.Sessions = sessions
	res.CodeOf(code)
	ctx.JSON(200, res)
}

func SendMessage(ctx *gin.Context) {
	var req SendMessageRequest
	var res SendMessageResponse
	var username string
	username = ctx.GetString("username")
	if username == "" {
		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
		return
	}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(200, res.CodeOf(common_code.CodeInvalidParams))
		return
	}
	content, code := service_aichat.SendMessage(req.ModelType, req.SessionID, req.Content, username)
	if code != common_code.CodeSuccess {
		ctx.JSON(200, res.CodeOf(code))
		return
	}
	res.Content = content
	res.CodeOf(code)
	ctx.JSON(200, res)
	return
}

func SendStreamMessage(ctx *gin.Context) {
	var req SendStreamMessageRequest
	var res SendStreamMessageResponse
	var username string
	// 验证请求数据是否合法
	username = ctx.GetString("username")
	if username == "" {
		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
		return
	}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(200, res.CodeOf(common_code.CodeInvalidParams))
	}
	// 开启流式响应先传输一个响应头，告诉客户端这是一个流式响应
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("X-Accel-Buffering", "no") // 禁止代理缓存
	// 流式发送
	flusher := ctx.Writer.(http.Flusher)
	// 回调函数，用于将模型生成的内容流式发送给客户端
	cb := func(msg string) {
		_, err := ctx.Writer.Write([]byte("data: " + msg + "\n\n"))
		if err != nil {
			log.Println("[SSE] Write error:", err)
			return
		}
		flusher.Flush() //  每次必须 flush
	}
	content, code := service_aichat.SendStreamMessage(req.ModelType, req.SessionID, req.Content, username, cb)
	if code != common_code.CodeSuccess {
		ctx.JSON(200, res.CodeOf(code))
		return
	}
	res.Content = content
	res.CodeOf(code)
	ctx.JSON(200, res)
}

func SessionDetail(ctx *gin.Context) {
	var req SessionDetailRequest
	var res SessionDetailResponse
	var username string
	username = ctx.GetString("username")
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(200, res.CodeOf(common_code.CodeInvalidParams))
		return
	}
	history, code := service_aichat.SessionDetail(username, req.SessionID)
	if code != common_code.CodeSuccess {
		ctx.JSON(200, res.CodeOf(code))
	}
	res.History = history
	res.CodeOf(code)
	ctx.JSON(200, res)
}

func CreateSession(ctx *gin.Context) {
	var username string
	// var req CreateSessionRequest
	var res CreateSessionResponse
	username = ctx.GetString("username")
	if username == "" {
		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
		return
	}
	sessionid, code := service_aichat.CreateSession(username)
	if code != common_code.CodeSuccess {
		ctx.JSON(200, res.CodeOf(code))
		return
	}
	res.SessionID = sessionid
	res.CodeOf(code)
	ctx.JSON(200, res)
}

// func CreateAndSendStreamSession(ctx *gin.Context) {
// 	var username string
// 	var req CreateAndSendSessionRequest
// 	var res CreateAndSendSessionResponse
// 	username = ctx.GetString("username")
// 	if username == "" {
// 		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
// 		return
// 	}
// 	err := ctx.BindJSON(&req)
// 	if err != nil {
// 		ctx.JSON(200, res.CodeOf(common_code.CodeInvalidParams))
// 	}
// 	content, code := service_aichat.CreateAndSendStreamSession(username, req.Content)
// 	if code != common_code.CodeSuccess {
// 		ctx.JSON(200, res.CodeOf(code))
// 		return
// 	}
// 	res.Content = content
// 	res.CodeOf(code)
// 	ctx.JSON(200, res)
// }

// func CreateAndSendSession(ctx *gin.Context) {
// 	var username string
// 	var req CreateAndSendSessionRequest
// 	var res CreateAndSendSessionResponse
// 	username = ctx.GetString("username")
// 	if username == "" {
// 		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
// 		return
// 	}
// 	err := ctx.BindJSON(&req)
// 	if err != nil {
// 		ctx.JSON(200, res.CodeOf(common_code.CodeInvalidParams))
// 	}
// 	content, code := service_aichat.CreateAndSendSession(username, req.Content)
// 	if code != common_code.CodeSuccess {
// 		ctx.JSON(200, res.CodeOf(code))
// 		return
// 	}
// 	res.Content = content
// 	res.CodeOf(code)
// 	ctx.JSON(200, res)
// }

// func DeleteSession(ctx *gin.Context) {
// 	var username string
// 	var req DeleteSessionRequest
// 	var res DeleteSessionResponse
// 	username = ctx.GetString("username")
// 	if username == "" {
// 		ctx.JSON(200, res.CodeOf(common_code.CodeUserNotExist))
// 		return
// 	}
// 	err := ctx.BindJSON(&req)
// 	if err != nil {
// 		ctx.JSON(200, res.CodeOf(common_code.CodeInvalidParams))
// 		return
// 	}
// 	code := service_aichat.DeleteSession(username, req.SessionID)
// 	if code != common_code.CodeSuccess {
// 		ctx.JSON(200, res.CodeOf(code))
// 		return
// 	}
// 	res.CodeOf(code)
// 	ctx.JSON(200, res)
// }
