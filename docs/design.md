- 依赖注入 FX
- WEB框架 GIN
- ORM框架 gorm
- 权限框架
- 


{{/*
{{ define "updateUser"}}
update sys_user
set {{$comma := false}}
{{if .Username}}{{$comma = false}}user_name = {{val .Username}}, {{end}}
{{if .LoginName}}{{if $comma}},{{end}}{{$comma = true}}login_name = {{val .LoginName}}{{end}}
{{if .Password}}{{if $comma}},{{end}}{{$comma = false}}password = {{val .Password}}{{end}}
{{if .LastLoginTime}}{{if $comma}},{{end}}{{$comma = true}}last_login_time = {{val .LastLoginTime}}{{end}}
where user_id = {{.UserId}}
{{end}}
*/}}