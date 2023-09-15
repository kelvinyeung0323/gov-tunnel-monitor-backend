package db

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ky/gov-tunnel-monitor-backend/base/config"
)

const TRANSACTION_CONTEXT_KEY = "transaction-context-key"

var mysqlConn *sqlx.DB

//TODO:加锁
func getDB() *sqlx.DB {
	// if mysqlConn == nil {
	// 	InitDB()
	// }
	return mysqlConn
}

func InitDB() {
	// dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
	// "root", "123456", "localhost", 3306, "go-tunnel")
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?allowNativePasswords=True&charset=utf8&parseTime=True&loc=UTC",
	// "root", "123456", "localhost", 33060, "go-tunnel")

	dsn := config.MySQLConfig.Username + ":" + config.MySQLConfig.Password + "@" + config.MySQLConfig.Url

	// sqlx连接池
	sqlxDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	sqlxDB.SetMaxOpenConns(500)
	sqlxDB.SetMaxIdleConns(100)
	mysqlConn = sqlxDB
}

//事务holder
//cnt 为嵌套计数，用于处理嵌套事务，这里没有考虑事务隔离级别和事务传播机制
type TransactionHolder struct {
	DB  *sqlx.DB
	Tx  *sqlx.Tx
	Cnt int
}

func BeginTx(ctx *gin.Context) {
	if ctx == nil {
		log.Println("没有ctx,不开启事务")
		return
	}
	//从Context中获取transaction holder,如果没有则新增一个
	//这里无需考需并发，一般一个请求只在一个routine中操作业务
	txHolder, exists := ctx.Get(TRANSACTION_CONTEXT_KEY)
	if !exists {
		txHolder = &TransactionHolder{
			DB: getDB(),
			Tx: nil,
		}
		ctx.Set(TRANSACTION_CONTEXT_KEY, txHolder)
	}

	th, ok := txHolder.(*TransactionHolder)
	if !ok {
		panic(fmt.Sprintf("error:获取事务管理错误:%v", txHolder))

	}

	if th.Tx == nil {
		var err error
		th.Tx, err = getDB().BeginTxx(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("error: 开启事务错误:%v\n", err))
		}
		th.Cnt++

	}

}

func CommitTx(ctx *gin.Context) {
	if ctx == nil {
		log.Println("没有ctx,不提交事务")
		return
	}
	txHolder, exists := ctx.Get(TRANSACTION_CONTEXT_KEY)
	if !exists {
		log.Println("error: 没有开启的事务.")
	}
	th, ok := txHolder.(*TransactionHolder)
	if !ok {
		panic(fmt.Sprintf("error:获取事务管理错误:%v\n", txHolder))
	}
	th.Cnt--
	if th.Cnt != 0 {
		return
	}
	err := th.Tx.Commit()
	if err != nil {

		panic(fmt.Sprintf("error: 提交事务错误:%v\n", err))
	}
}

type SQLCommon interface {
	sqlx.Ext
	Select(dest interface{}, query string, args ...interface{}) error
	// Get using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error
}

func GetConn(ctx *gin.Context) SQLCommon {
	if ctx == nil {
		return getDB()
	}

	txHolder, exists := ctx.Get(TRANSACTION_CONTEXT_KEY)
	if !exists {
		//没有开启事务，则用普通连接
		return getDB()
	}
	th, ok := txHolder.(*TransactionHolder)
	if !ok {
		panic(fmt.Sprintf("error:获取事务管理错误:%v\n", txHolder))
	}
	if th.Tx != nil {
		return th.Tx
	}
	return getDB()
}

//事务处理中间件
func TransactionMiddleware(ctx *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			txHolder, exists := ctx.Get(TRANSACTION_CONTEXT_KEY)
			if exists {
				th, ok := txHolder.(*TransactionHolder)
				if ok && th.Tx != nil {
					err := th.Tx.Rollback()
					if err != nil {
						log.Printf("error:%v", err)
					}

				}

			}
			panic(r)
		}

	}()

	ctx.Next()

}
