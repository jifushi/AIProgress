package aichat

import (
	"AiProgress/common/mysql"
	"AiProgress/model"
)

// sessions数据库操作

func CreateSession(session *model.Session) error {
	// 调用mysql的函数insertSession
	err := mysql.InsertSession(session)
	if err != nil {
		return err
	}
	return nil
}

// 会话信息缓存操作
