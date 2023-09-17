# 隧道监控系统
这个项目是为了学习go语言,所以一些常见的业务细节是没有做的，只对项目架构中的一些关注点进行了思考和代码编写；   
这个项目是之前是和一个小伙合作做的，因为后面才找到我，所以我负责前端，后端由小伙负责；  
项目基于开源的rouyi管理系统，然后加上监控的模块；最后因为学习go和blender、threejs所以把这个项目的前后端重写了一下；顺便作业个人作品展示
## 前端效果：
![image](https://github.com/kelvinyeung0323/tunnel-3d/blob/main/docs/pics/screen_home.png)  


## 项目架构设计
项目主要使用 go + gin + sqlx 

下面针对以下几项进行说明：
- 项目结构
- 权限验证
- 全局异常处理
- ORM
- 全局事务管理
- Websocket
- 监控线程

### 项目结构
参考java项目习惯将项目分层；  
routes对应java的controller,一些请求参数处理在route里做了，本来想把controller单独出来的，但发现这样做分层太多了，controller做的事本来就不多，所以跟route整合在一起；   
handler对应service层，负责业务逻辑的处理；  
repositories负责数据库操作； 
resources目录存放一些资源文件，使用embbed,打包时嵌入到go程序中；
其中，gin.Context贯穿reoutes->handler->respositories,用于一些全局处理；

### 权限验证
权限验证jwt,编写中间件对请求进行拦截验证权限；这里不多说，网上很多资料；

### 全局异常处理
一般java web应用都会设置一个全局异常处理类来对异常进行统一的处理；
所以项目里使用中间件进行统一异常处理；因为使用error的处理方式个人觉得过于啰唆；  
在repo里使用error的方式处理；因为repo里的逻辑没有太大的业务意义；
然后在handler中处理error,抛panic，这里的panic是有具体内业意义的；   
在中间件中recover panic;然后统一返回错误响应；


### ORM
个人不太喜欢gorm，喜欢java的mybatis的模式；  
本人对SQL比较熟悉，所以gorm的模式太别扭，不如写SQL来的方便；  
所以选择了sqlx+template的方式；   
将SQL写在模板中，通过template编译，这样就跟Mybatis差不多了；

### 全局事务处理
go不像java一样可以使用注解来配置事务处理;网上找的解决方案，个人觉得都不太好,过于复杂；在这个项里我把事务处理的逻辑提到中间件里进行处理；  
1.第一个请求到来，在中间件里开启事务，然后把数据库连接放在gin.Context里；    
2.repo层先从gin.Context中把连接合出来再进行操作；    
3.如果数据操作失败或业务逻辑错误，则抛出特定类型的panic;   
4.在中间件进行recover,然后对事务进行rollback;    
5.如果没有panic则commit;    


## Websocket
项目使用websocket跟前端进行实时交互


## 监控线程
监控线程负责定时轮询监控设备，把获取到的数据通过websocket发送到前端并保存相关日志；
通过socket连接设备，每种设备都有各自的交互协议；




