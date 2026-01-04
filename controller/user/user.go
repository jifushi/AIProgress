package user

import (
	"net/http"

	common_code "AiProgress/common/code"
	"AiProgress/controller"
	service_user "AiProgress/service/user"

	"github.com/gin-gonic/gin"
)

type (
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	LoginResponse struct {
		Token string `json:"token"`
		controller.Response
	}
	RegisterRequest struct {
		Password     string `json:"password"`
		Email        string `json:"email"`
		Verification string `json:"verification"`
	}
	RegisterResponse struct {
		Token string `json:"token"`
		controller.Response
	}
	VerificationEmailRequest struct {
		Email string `json:"email"`
	}
	VerificationEmailResponse struct {
		controller.Response
	}
)

func Login(ctx *gin.Context) {
	// 1. 接收请求参数
	req := new(LoginRequest)
	res := new(LoginResponse)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusOK, res.CodeOf(common_code.CodeInvalidParams))
		return
	}
	// 2. 校验参数,使用用户名或邮箱登录
	if (req.Username == "" && req.Email == "") || req.Password == "" {
		ctx.JSON(http.StatusOK, res.CodeOf(common_code.CodeInvalidParams))
		return
	}
	// 3. 去数据库对数据进行校验,使用用户名或邮箱登录
	token, code := service_user.Login(req.Username, req.Email, req.Password)
	if code != common_code.CodeSuccess {
		ctx.JSON(http.StatusOK, res.CodeOf(code))
		return
	}
	res.Token = token
	res.CodeOf(common_code.CodeSuccess)
	ctx.JSON(http.StatusOK, res)
}

func Register(ctx *gin.Context) {
	// 1. 接收请求参数
	req := new(RegisterRequest)
	res := new(RegisterResponse)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusOK, res.CodeOf(common_code.CodeInvalidParams))
		return
	}
	token, code := service_user.Register(req.Email, req.Password, req.Verification)
	if code != common_code.CodeSuccess {
		ctx.JSON(http.StatusOK, res.CodeOf(code))
		return
	}

	res.Token = token
	res.CodeOf(common_code.CodeSuccess)
	ctx.JSON(200, res)
}

// 邮箱验证，使用qq邮箱的smtp服务
func VerificationEmail(ctx *gin.Context) {
	// 1. 接收请求参数
	req := new(VerificationEmailRequest)
	res := new(VerificationEmailResponse)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 发送邮箱的逻辑
	code := service_user.SendVerificationEmail(req.Email)
	if code != common_code.CodeSuccess {
		ctx.JSON(http.StatusOK, res.CodeOf(code))
	}
	ctx.JSON(http.StatusOK, res.CodeOf(common_code.CodeSuccess))
}
