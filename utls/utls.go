package utls

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// 生成随机验证码
func CreateRandomNumber(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < n; i++ {
		code = fmt.Sprintf("%s%d", code, r.Intn(10))
	}
	return code
}

// MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
