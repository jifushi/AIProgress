package jwt

import (
	"log"
	"net/http"
	"strings"

	common_code "AiProgress/common/code"
	"AiProgress/controller"
	utls_jwt "AiProgress/utls/jwt"

	"github.com/gin-gonic/gin"
)

// jwt认证,拆解请求头中的token，
// 使用utls中的jwt包进行检验是否合法，
// 合法则返回token中的信息，否则返回错误
func Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var res *controller.Response
		var token string
		// 1. 获取请求头中的token
		authHeader := ctx.GetHeader("Authorization")
		// 拆解Authorization中的token,去除Bearer前缀
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// 如果没有token，去url中获取token
			token = ctx.Query("token")
		}

		if token == "" {
			ctx.JSON(http.StatusOK, res.CodeOf(common_code.CodeInvalidToken))
			// 中断请求
			// ctx.Abort() 中断请求，后面的处理函数不会执行
			ctx.Abort()
			return
		}

		log.Println("token is ", token)
		// 2. 校验token是否合法
		username, ok := utls_jwt.ParseToken(token)
		if !ok {
			ctx.JSON(http.StatusOK, res.CodeOf(common_code.CodeInvalidToken))
			// 中断请求
			// ctx.Abort() 中断请求，后面的处理函数不会执行
			ctx.Abort()
			return
		}

		// 3. 如果token合法，将token中的信息存入ctx中，方便后续使用
		ctx.Set("username", username)
		// 4. 继续执行后面的处理函数
		ctx.Next()
	}
}
