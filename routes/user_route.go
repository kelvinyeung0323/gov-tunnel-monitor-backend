package routes

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/web"
	"github.com/ky/gov-tunnel-monitor-backend/handlers"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

func registerUserRoute(r *gin.RouterGroup) {
	//查询用户
	r.GET("/users", func(ctx *gin.Context) {

		form := types.UserQueryForm{}
		if err := ctx.ShouldBindQuery(&form); err != nil {
			web.Err(web.VALID_ERROR)

		}
		users := handlers.QueryUser(ctx, &form)
		web.ReturnOK(ctx, &users)

	})

	//根据用户ID获取用户
	r.GET("/user/:userId", func(ctx *gin.Context) {
		s := ctx.Param("userId")
		userId, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("error:%v\n", err)
			web.Err(web.VALID_ERROR)
		}

		user := handlers.GetUserById(ctx, userId)
		web.ReturnOK(ctx, user)
	})

	//创建用户
	r.POST("/user", func(ctx *gin.Context) {
		form := &types.UserCreateForm{}
		if err := ctx.ShouldBind(&form); err != nil {
			log.Printf("error:%v", err)
			web.Err(web.VALID_ERROR)
		}

		user := handlers.CreateUser(ctx, form)
		web.ReturnOK(ctx, user)

	})

	//更新用户
	r.PUT("/user", func(ctx *gin.Context) {
		user := &types.User{}
		if err := ctx.ShouldBind(&user); err != nil {
			log.Printf("error:%v", err)
			web.Err(web.VALID_ERROR)
		}

		user = handlers.UpdateUser(ctx, user)
		web.ReturnOKWithMsg(ctx, user, "用户修改成功！")
	})

	//删除用户
	r.DELETE("/user/:userId", func(ctx *gin.Context) {

		s := ctx.Param("userId")
		userId, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("error:%v\n", err)
			web.Err(web.VALID_ERROR)
		}

		user := handlers.DelUser(ctx, userId)
		web.ReturnOKWithMsg(ctx, user, "用户已删除！")
	})
}
