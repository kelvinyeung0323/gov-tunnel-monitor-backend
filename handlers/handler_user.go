package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/db"
	"github.com/ky/gov-tunnel-monitor-backend/base/security"
	"github.com/ky/gov-tunnel-monitor-backend/base/web"
	repo "github.com/ky/gov-tunnel-monitor-backend/repositories"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

//查询用户
func QueryUser(ctx *gin.Context, queryForm *types.UserQueryForm) []types.User {
	return repo.QueryUser(ctx, queryForm)
}

//根据ID获取用户
func GetUserById(ctx *gin.Context, userId int) *types.User {
	user := repo.GetUserById(ctx, userId)
	if user == nil {
		web.BizErr("用户不存在")
	}
	return user
}

//创建用户
func CreateUser(ctx *gin.Context, form *types.UserCreateForm) *types.User {
	//开户事务
	db.BeginTx(ctx)
	//TODO:判断用户时间存在
	//TODO检查字段
	user := &types.User{}
	user.Username = form.Username
	user.LoginName = form.LoginName
	user.Password = security.EncryptPwd(form.Password)
	repo.CreateUser(ctx, user)

	//提交事务
	db.CommitTx(ctx)
	return user
}

//更新用户
//返回更新后的用户
func UpdateUser(ctx *gin.Context, user *types.User) *types.User {
	db.BeginTx(ctx)
	if user.UserId == nil {
		web.BizErr("用户不存在.")
	}

	oUser := repo.GetUserById(ctx, *user.UserId)
	if oUser == nil {
		web.BizErr("用户不存在.")
	}
	repo.UpdateUser(ctx, user)
	user = repo.GetUserById(ctx, *user.UserId)
	db.CommitTx(ctx)
	return user
}

//删除用户
func DelUser(ctx *gin.Context, userId int) *types.User {
	db.BeginTx(ctx)
	db.CommitTx(ctx)
	return nil
}

//删除多个用户
//返回删除后的用户
func DelUsers(ctx *gin.Context, userIds []int) []types.User {
	return nil
}
