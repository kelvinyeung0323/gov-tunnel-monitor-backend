package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/security"
	"github.com/ky/gov-tunnel-monitor-backend/base/web"
	repo "github.com/ky/gov-tunnel-monitor-backend/repositories"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

func HandleLogin(ctx *gin.Context) {
	var user types.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		web.Err(web.VALID_ERROR)
	}

	//TODO:从数据库读取用户信息
	// 校验用户名和密码是否正确
	userInDb, err := repo.GetUserByName(ctx, user.LoginName)
	if err != nil {
		log.Printf("查找用户失败%v", err)
		web.Err(web.LOGIN_UNKNOWN)
	}

	//TODO:加码码码

	if *userInDb.Password == *security.EncryptPwd(user.Password) {
		// 生成Token
		tokenString, _ := security.GenToken(*user.UserId, *user.Username)
		web.ReturnOK(ctx, gin.H{"token": tokenString})
		return
	}
	web.Err(web.LOGIN_ERROR)

}
