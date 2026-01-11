package rabbitmq

import (
	"encoding/json"

	dao_aichat "AiProgress/dao/aichat"
	"AiProgress/model"

	"github.com/streadway/amqp"
)

type MessageMQParam struct {
	SessionID string `json:"session_id"`
	Content   string `json:"content"`
	UserName  string `json:"user_name"`
	IsUser    bool   `json:"is_user"`
}

// 生产者调用程序，将数据转化为json格式，方便后续插入消息队列中
func GenerateMessageMQParam(sessionid string, content string, username string, isuser bool) []byte {
	param := MessageMQParam{
		SessionID: sessionid,
		Content:   content,
		UserName:  username,
		IsUser:    isuser,
	}
	data, _ := json.Marshal(param)
	return data
}

// 消费者调用程序，将消息队列中的数据转化为结构体，方便后续插入数据库中
func MQMessage(msg *amqp.Delivery) error {
	var param MessageMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	newMsg := &model.Message{
		SessionID: param.SessionID,
		Content:   param.Content,
		UserName:  param.UserName,
		IsUser:    param.IsUser,
	}
	// 消费者异步插入到数据库中
	dao_aichat.CreateMessage(newMsg)
	return nil
}
