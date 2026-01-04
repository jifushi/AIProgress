package controller

import "AiProgress/common/code"

// 为controller提供公共功能
// 例如：请求参数的校验，错误处理等，统一管理响应的格式，统一管理错误码等。
type Response struct {
	StatusCode code.Code `json:"status_code"`
	StatusMsg  string    `json:"status_msg"`
}

func (r *Response) CodeOf(code code.Code) Response {
	if nil == r {
		r = new(Response)
	}
	r.StatusCode = code
	r.StatusMsg = code.Msg()
	return *r
}

func (r *Response) Success() {
	r.CodeOf(code.CodeSuccess)
}
