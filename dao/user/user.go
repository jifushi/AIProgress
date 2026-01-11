package user

import (
	"AiProgress/common/mysql"
	"AiProgress/model"
)

// func IsExistUser(email string) (bool, *model.User) {
// 	// 查询数据库判断用户是否存在
// 	var user model.User
// 	var err error
// 	user, err = mysql.GetUserByEmail(email)
// 	return false, &user
// }

func InsertUser(username string, password string, email string) (*model.User, error) {
	// 插入数据库
	var user model.User
	user = model.User{
		Name:     username,
		UserName: username,
		Password: password,
		Email:    email,
	}
	err := mysql.InsertUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*model.User, error) {
	// 查询数据库
	var user model.User
	err := mysql.GetUserByUsername(username, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	// 查询数据库
	var user model.User
	err := mysql.GetUserByEmail(email, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
