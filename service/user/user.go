package user

import (
	common_code "AiProgress/common/code"
	common_email "AiProgress/common/email"
	common_redis "AiProgress/common/redis"
	dao_user "AiProgress/dao/user"
	"AiProgress/model"
	"AiProgress/utils"
	utils_jwt "AiProgress/utils/jwt"
)

func Login(username, email, password string) (string, common_code.Code) {
	// 1. 检查用户登录方式是使用用户名还是邮箱登录
	var userInformation *model.User
	var err error
	password = utils.MD5(password)
	if username != "" {
		// 2. 检查用户名是否存在
		userInformation, err = dao_user.GetUserByUsername(username)
		if err != nil {
			return "", common_code.CodeUserNotExist
		}
		// 3. 检查密码是否正确
		if userInformation.Password != password {
			return "", common_code.CodePasswordError
		}
		// 4. 生成token
		token, err := utils_jwt.Create(userInformation.ID, userInformation.UserName)
		if err != nil {
			return "", common_code.CodeServerBusy
		}
		return token, common_code.CodeSuccess
	} else {
		// 2. 检查邮箱是否存在
		userInformation, err = dao_user.GetUserByEmail(email)
		if err != nil {
			return "", common_code.CodeEmailNotExist
		}
		// 3. 检查密码是否正确
		if userInformation.Password != password {
			return "", common_code.CodePasswordError
		}
		// 4. 生成token
		token, err := utils_jwt.Create(userInformation.ID, userInformation.UserName)
		if err != nil {
			return "", common_code.CodeServerBusy
		}
		return token, common_code.CodeSuccess
	}
}

func Register(email, password, verification string) (string, common_code.Code) {
	// 1. 检查邮箱是否已经注册
	var userInformation *model.User
	var err error
	if userInformation, err = dao_user.GetUserByEmail(email); err == nil {
		return "", common_code.CodeUserExist
	}
	// 2. 检查验证码是否正确
	if ok, _ := common_redis.CheckVerifyCodeEmail(email, verification); !ok {
		return "", common_code.CodeInvalidCaptcha
	}
	// 3. 生成用户名
	username := utils.CreateRandomNumber(11)
	password = utils.MD5(password)
	// 4. 添加用户到数据库
	userInformation, err = dao_user.InsertUser(username, password, email)
	if err != nil {
		return "", common_code.CodeServerBusy
	}
	// 5.将生成的账号发送到用户邮箱
	err = common_email.SendUsername(email, username)
	if err != nil {
		return "", common_code.CodeServerBusy
	}
	// 6. 生成token
	token, err := utils_jwt.Create(userInformation.ID, userInformation.UserName)
	// 7. 返回token
	if err != nil {
		return "", common_code.CodeServerBusy
	}
	return token, common_code.CodeSuccess
}

func SendVerificationEmail(email string) common_code.Code {
	// 1. 检查邮箱是否已经注册
	verification := utils.CreateRandomNumber(6)
	// 将验证码存入redis中，设置过期时间为2分钟
	err := common_redis.SetVerifyCodeEmail(email, verification)
	if err != nil {
		return common_code.CodeServerBusy
	}
	err = common_email.SendVerifyCode(email, verification)
	if err != nil {
		return common_code.CodeServerBusy
	}
	return common_code.CodeSuccess
}
