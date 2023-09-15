package repo

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/db"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

func GetUserById(ctx *gin.Context, id int) *types.User {
	u := &types.User{}
	sql := db.GetSqlString("getUserById", id)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(ctx).Get(u, sql)
	if err != nil {
		panic(err)
	}
	return u
}

func GetUserByName(ctx *gin.Context, username *string) (*types.User, error) {
	u := &types.User{}
	sql := db.GetSqlString("getUserByName", username)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(ctx).Get(u, sql)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func QueryUser(ctx *gin.Context, queryForm *types.UserQueryForm) []types.User {

	users := []types.User{}
	sql := db.GetSqlString("queryUsers", queryForm)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(ctx).Select(&users, sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		panic(err)
	}

	return users
}

func CreateUser(ctx *gin.Context, user *types.User) {

	sql := db.GetSqlString("createUser", user)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		panic(err)
	}

	if lastId, e := r.LastInsertId(); e == nil {
		userId := int(lastId)
		user.UserId = &userId
	} else {
		panic(e)
	}

}

func UpdateUser(ctx *gin.Context, user *types.User) {
	sql := db.GetSqlString("updateUser", user)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		panic(err)
	}
	cnt, _ := r.RowsAffected()
	log.Printf("Rows Affected: %v", cnt)
}

func ChangeUserPwd(ctx *gin.Context, userId *int, password *string) {
	user := &types.User{}
	user.UserId = userId
	user.Password = password
	sql := db.GetSqlString("changeUserPasswd", user)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		panic(err)
	}
	cnt, _ := r.RowsAffected()
	log.Printf("Rows Affected: %v", cnt)
}

func UdateLoginTime(ctx *gin.Context, userId *int, lastLoginTime *types.Time) {
	user := &types.User{}
	user.UserId = userId
	user.LastLoginTime = lastLoginTime
	sql := db.GetSqlString("updateLoginTime", user)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		panic(err)
	}
	cnt, _ := r.RowsAffected()
	log.Printf("Rows Affected: %v", cnt)
}
