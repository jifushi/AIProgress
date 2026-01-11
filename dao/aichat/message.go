package aichat

import (
	"log"

	"AiProgress/common/mysql"
	"AiProgress/common/redis"
	"AiProgress/model"
)

/*
***********************
对mysql的操作
***********************
*/
func CreateMessage(msg *model.Message) {
	err := mysql.InsertMessage(msg)
	if err != nil {
		log.Fatal(err)
	}
	// return msg, err
}

func MySQLGetHistoryMessage(username, sessionid string) ([]*model.History, error) {
	messages, err := mysql.SelectMessages(username, sessionid)
	if err != nil {
		return nil, err
	}
	history := make([]*model.History, len(messages))
	for i := range messages {
		history[i] = &model.History{
			IsUser:  messages[i].IsUser,
			Content: messages[i].Content,
		}
	}
	return history, nil
}

func GetSessions(username string) ([]model.SessionInfo, error) {
	sessions, err := mysql.SelectSessions(username)
	if err != nil {
		return nil, err
	}
	sessioninfo := make([]model.SessionInfo, len(sessions))
	for i := range sessions {
		sessioninfo[i] = model.SessionInfo{
			SessionID: sessions[i].ID,
			Title:     sessions[i].Title,
		}
	}
	return sessioninfo, nil
}

/*
***********************
对redis的操作
***********************
*/
func RedisGetHistoryMessage(username, sessionid string) ([]*model.History, error) {
	history, err := redis.HGetMessages(username, sessionid)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// 初始化，将历史上下文放入redis做缓存
func RedisSetHistoryMessage(username string, sessionid string, history []*model.History) error {
	return redis.HSetMessages(username, sessionid, history)
}

// 添加对话信息到redis中
func RedisAddHistoryMessage(username string, sessionid string, oneChatMessage []*model.History) error {
	return redis.HUpdateMessages(username, sessionid, oneChatMessage)
}
