package main

import (
	"fmt"
	"log"

	"AiProgress/common/mysql"
	"AiProgress/common/rabbitmq"
	"AiProgress/common/redis"
	"AiProgress/config"
	"AiProgress/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func server() {
	var r *gin.Engine
	r = routers.RouteController()
	r.Run("127.0.0.1:8080")
}

func main() {
	godotenv.Load() // 加载环境变量文件
	config.LoadConfig()
	err := mysql.Init()
	if err != nil {
		fmt.Println("数据库初始化失败:", err)
	}
	redis.Init()
	log.Println("redis init success  ")
	rabbitmq.Init()
	log.Println("rabbitmq init success  ")
	server()
}
