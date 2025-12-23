package routers

import (
	"AiProgress/controller/user"

	"github.com/gin-gonic/gin"
)

func userRouter(rg *gin.RouterGroup) {
	rg.POST("/login", user.Login)
	rg.POST("/register", user.Register)
	verificationGroup := rg.Group("/verification")
	verificationGroup.POST("/email", user.VerificationEmail)
	verificationGroup.POST("/phone", user.VerificationPhone)
}
