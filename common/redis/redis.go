package redis

import (
	"context"
	"fmt"
	"time"

	"AiProgress/config"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

func Init() {
	// 初始化 Redis 客户端
	addr := fmt.Sprintf("%s:%d", config.GetRedisConfig().Host, config.GetRedisConfig().Port)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr, // 替换为你的 Redis 地址和端口
		Password: "",   // 如果有密码，替换为密码
		DB:       0,    // 使用默认数据库
	})
}

func SetVerifyCodeEmail(email string, code string) error {
	expire := 2 * time.Minute
	return Rdb.Set(ctx, "verify_email_"+email, code, expire).Err()
}

func CheckVerifyCodeEmail(email string, code string) (bool, error) {
	verifyCode, err := Rdb.Get(ctx, "verify_email_"+email).Result()
	if err != nil {
		return false, err
	}
	if verifyCode == code {
		// 删除验证码
		Rdb.Del(ctx, "verify_email_"+email)
		return true, nil
	}
	return false, nil
}
