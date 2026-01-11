package mysql

import (
	"fmt"
	"log"
	"time"

	"AiProgress/config"
	"AiProgress/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() error {
	mysqlConfig := config.GetMysqlConfig()
	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User, // 使用第一个用户名
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database)
	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	// 连接数据库
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), gormConfig)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
		return err
	}

	log.Println("数据库连接成功")
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 迁移数据库模型
	return Migrate()
}

func Migrate() error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Message{},
		&model.Session{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
		return err
	}
	log.Println("数据库迁移成功")
	return nil
}

/*
*********************************

	对用户做增删改查操作

**********************************
*/
func InsertUser(user *model.User) error {
	err := db.Create(user).Error
	return err
}

func GetUserByUsername(username string, user *model.User) error {
	err := db.Where("user_name = ?", username).First(user).Error
	return err
}

func GetUserByEmail(email string, user *model.User) error {
	err := db.Where("email =?", email).First(user).Error
	return err
}

/*
*********************************

	对消息做增删改查操作

**********************************
*/
func InsertMessage(message *model.Message) error {
	err := db.Create(message).Error
	return err
}

func SelectMessages(username, sessionid string) ([]model.Message, error) {
	var messages []model.Message
	err := db.Table("messages").Where("user_name =? AND session_id =?", username, sessionid).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

/*
*********************************

	对会话做增删改查操作

**********************************
*/
func InsertSession(session *model.Session) error {
	err := db.Create(session).Error
	return err
}

func SelectSessions(username string) ([]model.Session, error) {
	var sessions []model.Session
	err := db.Table("sessions").Select("DISTINCT id").Where("user_name = ?", username).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
