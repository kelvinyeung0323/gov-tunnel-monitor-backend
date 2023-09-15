package types

/**
* 用户
**/
type User struct {
	UserId        *int    `json:"userId" db:"user_id"`
	LoginName     *string `json:"loginName" db:"login_name"`
	Username      *string `json:"username" db:"user_name"`
	Password      *string `json:"password,omitempty" db:"password"`
	LastLoginTime *Time   `json:"lastLoginTime" db:"last_login_time" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedAt     *Time   `json:"-" db:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	CreatedAt     *Time   `json:"-" db:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
}

//查询表单
type UserQueryForm struct {
	Username *string `form:"username"`
}

//创建用户表单
type UserCreateForm struct {
	LoginName *string `json:"loginName" db:"login_name" binding:"required"`
	Username  *string `json:"username" db:"user_name" binding:"required"`
	Password  *string `json:"password,omitempty" db:"password" binding:"required"`
}

/**
* 角色
**/
type Role struct {
	RoleId    int
	RoleName  string
	RoleCode  string
	UpdatedAt NullTime
	CreatedAt NullTime
}

/**
*  资源
**/
type Resource struct {
	ResId     int
	ResCode   string
	ResName   string
	UpdatedAt NullTime
	CreatedAt NullTime
}
