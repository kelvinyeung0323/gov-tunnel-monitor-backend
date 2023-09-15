package web

import (
	"log"

	"github.com/gin-gonic/gin"
)

type webError struct {
	Code int64
	Msg  string
}

var (
	//错误代码
	LOGIN_UNKNOWN = &webError{202, "用户不存在"}
	LOGIN_ERROR   = &webError{203, "账号或密码错误"}
	TOKEN_ERROR   = &webError{204, "权限验证错误"}
	VALID_ERROR   = &webError{300, "参数错误"}
	ERROR         = &webError{400, "操作失败"}
	UNAUTHORIZED  = &webError{401, "您还未登录"}
	NOT_FOUND     = &webError{404, "资源不存在"}
	INNER_ERROR   = &webError{500, "系统发生异常"}
	BIZ_FAIL      = &webError{501, "业务错误"}
)

func Err(err *webError) {
	panic(err)
}

//其他业务相关的错误统一使用此方法
func BizErr(msg string) {
	panic(&webError{Code: 501, Msg: msg})
}

//统一异常处理中间件
func ErrorHandleMiddleware(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case *webError:
				log.Printf("panic: %v\n", t.Msg)
				ReturnFail(c, t)
			default:
				log.Printf("panic: %v\n", t)
				ReturnFail(c, INNER_ERROR)
			}
			c.Abort()
		}
	}()
	c.Next()
}
