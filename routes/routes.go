package routes

import (
	"fmt"
	"reflect"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/security"
	"github.com/ky/gov-tunnel-monitor-backend/base/websocket"
	"github.com/ky/gov-tunnel-monitor-backend/handlers"
)

func InitRoute(r *gin.Engine) error {

	//跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://foo.com"},                        //允许跨域发来请求的网站
		AllowMethods:  []string{"GET", "POST", "DELETE", "PUT", "OPTION"}, //允许请求的方法
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
	}))

	//websocket
	hub := websocket.NewHub()
	go hub.Run()
	r.Any("/ws", func(ctx *gin.Context) {
		websocket.WebsocketHandler(hub, ctx)
	})
	r.POST("/login", handlers.HandleLogin)
	//路由
	api := r.Group("/api/", security.JWTAuthMiddleware())
	{
		registerUserRoute(api)
		registerDeviceRoute(api)
	}

	return nil
}

/**
* 预处理函数入参，这样无需每个handler都以*gin.context为参数，写上需要的入参，将自动注入需要参数；
* TODO: GO 不支持反射获取函数参数的名称，暂不做
 */
func preHandle(fn interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		fmt.Println("注册的handler不是func类型")
	}

	return func(ctx *gin.Context) {
		//TODO:预处理入参
		fmt.Println("预处理入参")
		funcVal := reflect.ValueOf(fn)
		funcType := reflect.TypeOf(fn)

		paramList := []reflect.Value{}

		//自动转换类型
		for i := 0; i < funcType.NumIn(); i++ {
			paramType := funcType.In(i)
			funcVal.Type().Name()

			paramName := paramType.String()

			fmt.Printf("--cnt:%d,==num:%d ---%v\n", funcType.NumIn(), i, paramName)

			//gin包内的对像
			switch paramName {
			case "*gin.Context":
				paramList = append(paramList, reflect.ValueOf(ctx))
				fmt.Println("gin.Context end......")
			case "*http.Request":
				paramList = append(paramList, reflect.ValueOf(ctx.Request))
			case "*http.Response":
				paramList = append(paramList, reflect.ValueOf(ctx.Request.Response))
			default:
				fmt.Printf("====default:==%v", paramType.Kind())
				if paramType.Kind() == reflect.Struct ||
					(paramType.Kind() == reflect.Pointer &&
						paramType.Elem().Kind() == reflect.Struct) {

					fmt.Println("-- matched --" + paramName)
					//是否pointer类型
					var paramVal reflect.Value
					if paramType.Kind() == reflect.Pointer {
						fmt.Printf("it is a pointer\n")
						paramVal = reflect.New(paramType.Elem())
						if ctx.ShouldBind(paramVal.Interface()) == nil {
							fmt.Printf("---should bind --- %v\n", paramVal)
							paramList = append(paramList, paramVal)
						} else {
							//参数类型不对

							paramList = append(paramList, paramVal)
						}

					} else {
						fmt.Printf("it is not a pointer :%v\n", paramType)
						paramVal = reflect.New(paramType)
						fmt.Printf("value is :%v \n", paramVal)
						if ctx.ShouldBind(paramVal.Interface()) == nil {
							fmt.Printf("---should bind --- %v\n", paramVal)
							paramList = append(paramList, paramVal.Elem())
						} else {
							//参数类型不对
							paramList = append(paramList, paramVal.Elem())
						}
					}

				} else {
					//普通类型
					fmt.Println("-- 普通类型--")
					pName := paramType.Name()
					fmt.Printf("pName:%v\n", pName)
					pVal, _ := ctx.Get(pName)
					fmt.Printf("pVal:%v\n", pVal)
					paramList = append(paramList, reflect.ValueOf(pVal))

				}

			}
		}
		fmt.Println("---call-----" + funcVal.String())
		fmt.Printf("---call %v\n", paramList)
		r := funcVal.Call(paramList)
		if len(r) > 0 {
			fmt.Printf("---r is %v\n", r[0])
		} else {
			fmt.Printf("----no result\n")
		}

	}
}
