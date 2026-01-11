package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"AiProgress/model"

	"github.com/cloudwego/eino/schema"
)

// 生成随机验证码
func CreateRandomNumber(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < n; i++ {
		code = fmt.Sprintf("%s%d", code, r.Intn(10))
	}
	return code
}

// MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 数据转换为模型使用的消息格式
func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsags := make([]*schema.Message, 0, len(msgs))
	for _, msg := range msgs {
		role := schema.Assistant
		if msg.IsUser {
			role = schema.User
		}
		schemaMsags = append(schemaMsags, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}
	return schemaMsags
}

// 模型数据转换为本地使用的消息格式
func ConvertToModelMessage(sessionid string, username string, schemaMsag *schema.Message) *model.Message {
	var modelMsg *model.Message
	isUser := false
	if schemaMsag.Role == schema.User {
		isUser = true
	}
	modelMsg = &model.Message{
		SessionID: sessionid,
		UserName:  username,
		Content:   schemaMsag.Content,
		IsUser:    isUser,
	}

	return modelMsg
}
