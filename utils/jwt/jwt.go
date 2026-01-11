package jwt

import (
	"time"

	"AiProgress/config"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 需要两个函数，
// jwt生成
func Create(id int64, username string) (string, error) {
	// 创建的jwt需要考虑，
	// 1. 过期时间
	// 2. 加密方式
	// 3. 加密密钥
	// 4. 加密内容
	claims := Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.GetJwtConfig().Expiration))), // 过期时间
			Issuer:    config.GetJwtConfig().Issuer,                                                                      // 加密方式
			Subject:   config.GetJwtConfig().Subject,                                                                     // 主题
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                    // 签发时间
		},
	}
	// 5. 加密后的token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJwtConfig().Secret))
}

// jwt解析，用于判断token是否合法
func ParseToken(token string) (string, bool) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetJwtConfig().Secret), nil
	})
	if !t.Valid || err != nil {
		return "", false
	}
	return claims.Username, true
}
