package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构体
type Config struct {
	Mysql    MysqlConfig    `yaml:"mysql"`
	Redis    RedisConfig    `yaml:"redis"`
	Rabbitmq RabbitmqConfig `yaml:"rabbitmq"`
	Jwt      JwtConfig      `yaml:"jwt"`
	Aliyun   AliyunConfig   `yaml:"aliyun"`
	QqStmp   QqStmpConfig   `yaml:"qq_stmp"`
}

// MysqlConfig MySQL配置结构体
type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// RabbitmqConfig RabbitMQ配置结构体
type RabbitmqConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"vhost"`
}

// JwtConfig JWT配置结构体
type JwtConfig struct {
	Expiration int    `yaml:"expiration"`
	Issuer     string `yaml:"issuer"`
	Subject    string `yaml:"subject"`
	Secret     string `yaml:"secret"`
}

// AliyunConfig 阿里云配置结构体
type AliyunConfig struct {
	Phone  AliyunPhoneConfig  `yaml:"phone"`
	Qianyi AliyunQianyiConfig `yaml:"qianyi"`
}

// AliyunPhoneConfig 阿里云短信配置结构体
type AliyunPhoneConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

// AliyunQianyiConfig 阿里云千义配置结构体
type AliyunQianyiConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

// QqStmpConfig QQ SMTP配置结构体
type QqStmpConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// GlobalConfig 全局配置变量
var GlobalConfig Config

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 获取配置文件路径
	configPath := "config/config.yaml"

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML配置
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	log.Println("配置文件加载成功")
	return nil
}

// GetMysqlConfig 获取MySQL配置
func GetMysqlConfig() MysqlConfig {
	return GlobalConfig.Mysql
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig() RedisConfig {
	return GlobalConfig.Redis
}

// GetRabbitmqConfig 获取RabbitMQ配置
func GetRabbitmqConfig() RabbitmqConfig {
	return GlobalConfig.Rabbitmq
}

// GetJwtConfig 获取JWT配置
func GetJwtConfig() JwtConfig {
	return GlobalConfig.Jwt
}

// GetAliyunConfig 获取阿里云配置
func GetAliyunConfig() AliyunConfig {
	return GlobalConfig.Aliyun
}

// GetQqStmpConfig 获取QQ SMTP配置
func GetQqStmpConfig() QqStmpConfig {
	return GlobalConfig.QqStmp
}
