package aichat

import (
	"context"
	"log"

	common_aichat "AiProgress/common/aichat"
	common_code "AiProgress/common/code"
	dao_aichat "AiProgress/dao/aichat"

	"AiProgress/model"

	"github.com/google/uuid"
)

var ctx = context.Background()

func GetSessions(username string) ([]model.SessionInfo, common_code.Code) {
	// 直接从mysql中获取所有会话信息
	// 从数据库中获取会话信息
	var sessions []model.SessionInfo
	sessions, err := dao_aichat.GetSessions(username)
	if err != nil {
		return nil, common_code.CodeSQLError
	}
	return sessions, common_code.CodeSuccess
}

func SendMessage(mode_type, sessionid, content, username string) (string, common_code.Code) {
	// 1.先通过mode_type生成模型结构体
	var aihelper *common_aichat.AIHelper
	var err error
	ModeConfig := map[string]interface{}{
		"baseURL":   "http://127.0.0.1/myself.com",
		"modelName": "llama3.1",
	}
	aihelper, err = common_aichat.NewAIChat(ctx, username, sessionid, mode_type, ModeConfig)
	// 2.通过redis获取历史消息，如果没有再从数据库中获取
	// 查询当前会话窗口的上下文
	var historymessage []*model.History
	historymessage, err = CheckMessage(username, sessionid)
	// 3.向模型结构体的[]model.message中添加上下文内容
	aihelper.InitAIChatMessage(username, sessionid, historymessage)
	// 4.调用模型生成回复
	var airesponse *model.Message
	airesponse, err = aihelper.GenerateResponse(username, ctx, content)
	if err != nil {
		return "", common_code.AIModelFail
	}
	// 5.将一次对话的内容存入redis中
	OneChatMessage := []*model.History{
		{
			IsUser:  true,
			Content: content,
		},
		{
			IsUser:  false,
			Content: airesponse.Content,
		},
	}
	err = dao_aichat.RedisAddHistoryMessage(username, sessionid, OneChatMessage)
	if err != nil {
		return "", common_code.CodeRedisError
	}
	return airesponse.Content, common_code.CodeSuccess
}

func SendStreamMessage(mode_type, sessionid, content, username string, cb common_aichat.StreamCallback) (string, common_code.Code) {
	// 1.先通过mode_type生成模型结构体
	var aihelper *common_aichat.AIHelper
	var err error
	ModeConfig := map[string]interface{}{
		"baseURL":   "http://127.0.0.1/myself.com",
		"modelName": "llama3.1",
	}
	aihelper, err = common_aichat.NewAIChat(ctx, username, sessionid, mode_type, ModeConfig)
	// 2.通过redis获取历史消息，如果没有再从数据库中获取
	// 查询当前会话窗口的上下文
	var historymessage []*model.History
	// 先从redis中获取，如果没有再从数据库中获取
	historymessage, err = CheckMessage(username, sessionid)
	if err != nil {
		return "", common_code.CodeRedisError
	}
	// 3.向模型结构体的[]model.message中添加上下文内容
	aihelper.InitAIChatMessage(username, sessionid, historymessage)
	// 4.调用模型生成回复
	// 接收完成的ai响应放入redis缓存中
	var airesponse *model.Message
	airesponse, err = aihelper.StreamResponse(username, ctx, cb, content)
	if err != nil {
		return "", common_code.AIModelFail
	}
	// 5.将一次对话的内容存入redis中
	OneChatMessage := []*model.History{
		{
			IsUser:  true,
			Content: content,
		},
		{
			IsUser:  false,
			Content: airesponse.Content,
		},
	}
	err = dao_aichat.RedisAddHistoryMessage(username, sessionid, OneChatMessage)
	if err != nil {
		return "", common_code.CodeRedisError
	}
	return "1", common_code.CodeSuccess
}

func SessionDetail(username, sessionid string) ([]*model.History, common_code.Code) {
	// 先从redis中获取，如果没有再从数据库中获取
	var historymessage []*model.History
	historymessage, err := dao_aichat.RedisGetHistoryMessage(username, sessionid)
	if err != nil {
		historymessage, err = dao_aichat.MySQLGetHistoryMessage(username, sessionid)
		if err != nil {
			return nil, common_code.CodeSQLError
		}
		// 将历史消息存入redis中
		err = dao_aichat.RedisSetHistoryMessage(username, sessionid, historymessage)
		if err != nil {
			return nil, common_code.CodeRedisError
		}
	}
	return historymessage, common_code.CodeSuccess
}

// 创建新的会话，返回会话id
func CreateSession(username string) (string, common_code.Code) {
	// 生成一个uuid作为sessionid
	newSession := &model.Session{
		ID:       uuid.New().String(),
		UserName: username,
		Title:    uuid.New().String(), // title暂时用uuid作为默认值
	}
	err := dao_aichat.CreateSession(newSession)
	if err != nil {
		log.Println("CreateStreamSessionOnly CreateSession error:", err)
		return "", common_code.CodeServerBusy
	}
	return newSession.ID, common_code.CodeSuccess
}

func CreateAndSendStreamSession(username, content string) (string, common_code.Code) {
	return "1", common_code.CodeSuccess
}

func CreateAndSendSession(username, content string) (string, common_code.Code) {
	return "1", common_code.CodeSuccess
}

func DeleteSession(username, sessionid string) common_code.Code {
	return common_code.CodeSuccess
}

// 将重复使用的查询用户会话消息的函数封装为一个函数
func CheckMessage(username string, sessionid string) ([]*model.History, error) {
	// 2.通过redis获取历史消息，如果没有再从数据库中获取
	// 查询当前会话窗口的上下文
	var historymessage []*model.History
	// 先从redis中获取，如果没有再从数据库中获取
	historymessage, err := dao_aichat.RedisGetHistoryMessage(username, sessionid)
	if err != nil {
		historymessage, err = dao_aichat.MySQLGetHistoryMessage(username, sessionid)
		if err != nil {
			return nil, err
		}
		// 将历史消息存入redis中
		err = dao_aichat.RedisSetHistoryMessage(username, sessionid, historymessage)
	}
	return historymessage, nil
}
