package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"AiProgress/config"
	"AiProgress/model"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func Init() {
	// 初始化 Redis 客户端
	addr := fmt.Sprintf("%s:%d", config.GetRedisConfig().Host, config.GetRedisConfig().Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr, // 替换为你的 Redis 地址和端口
		Password: "",   // 如果有密码，替换为密码
		DB:       0,    // 使用默认数据库
	})
}

func SetVerifyCodeEmail(email string, code string) error {
	expire := 2 * time.Minute
	return rdb.Set(ctx, "verify_email_"+email, code, expire).Err()
}

func CheckVerifyCodeEmail(email string, code string) (bool, error) {
	verifyCode, err := rdb.Get(ctx, "verify_email_"+email).Result()
	if err != nil {
		return false, err
	}
	if verifyCode == code {
		// 删除验证码
		rdb.Del(ctx, "verify_email_"+email)
		return true, nil
	}
	return false, nil
}

/*
************************

	对会话信息做缓存操作
************************
*/
//用于启动服务时存储用户的会话信息（从mysql中获取数据）
// 存储用户的会话信息
func HSetMessages(username, sessionid string, messages []*model.History) error {
	userKey := fmt.Sprintf("user:%s", username)
	// 序列化消息
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}
	// 存储会话信息到 Redis Hash 中
	err = rdb.HSet(ctx, userKey, sessionid, data).Err()
	if err != nil {
		return fmt.Errorf("set session failed: %v", err)
	}
	// 设置过期时间并处理错误
	if err := rdb.Expire(ctx, userKey, 30*time.Minute).Err(); err != nil {
		return fmt.Errorf("set expire failed: %v", err)
	}
	return nil
}

// 用于用户发送消息时更新用户的会话信息（从redis中获取数据）
func HUpdateMessages(username, sessionid string, messages []*model.History) error {
	userKey := fmt.Sprintf("user:%s", username)
	// 检查用户是否存在
	// 获取会话信息
	data, err := rdb.HGet(ctx, userKey, sessionid).Result()
	// 证明该用户可能是新用户没有记录过，或者是新创建的会话没有记录，直接插入值
	if err == redis.Nil || data == "" {
		err = HSetMessages(username, sessionid, messages)
		if err != nil {
			return fmt.Errorf("set session failed: %v", err)
		}
	}
	// 反序列化消息
	var oldMessages []*model.History
	err = json.Unmarshal([]byte(data), &oldMessages)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %v", err)
	}
	// 合并消息
	oldMessages = append(oldMessages, messages...)
	// 序列化消息
	newdata, err := json.Marshal(oldMessages)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}
	// 存储会话信息到 Redis Hash 中
	err = rdb.HSet(ctx, userKey, sessionid, newdata).Err()
	if err != nil {
		return fmt.Errorf("set session failed: %v", err)
	}
	return nil
}

// 获取指定用户的历史会话信息
func HGetMessages(username, sessionid string) ([]*model.History, error) {
	userKey := fmt.Sprintf("user:%s", username)
	// 检查用户是否存在
	exists, err := rdb.Exists(ctx, userKey).Result()
	if err != nil {
		return nil, fmt.Errorf("exists check failed: %v", err)
	}
	if exists == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// 获取会话信息
	data, err := rdb.HGet(ctx, userKey, sessionid).Result()
	if err != nil {
		return nil, fmt.Errorf("get session failed: %v", err)
	}
	// 反序列化消息
	var messages []*model.History
	err = json.Unmarshal([]byte(data), &messages)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}
	return messages, nil
}
