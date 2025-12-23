package jwt

import (
	"github.com/gin-gonic/gin"
)

// 需要两个函数，
// jwt生成
func Create() {
}

// jwt认证
func Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	}
}
