package security

import (
	"crypto/md5"
	"fmt"
)

//MD5加密码
func EncryptPwd(pwd *string) *string {
	data := []byte(*pwd)
	hash := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", hash) //将[]byte转成16进制
	return &md5str1
}
