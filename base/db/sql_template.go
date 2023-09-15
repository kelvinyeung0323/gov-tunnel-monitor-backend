package db

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/ky/gov-tunnel-monitor-backend/types"
)

var sqlTemplate *template.Template

//解释模板函数
//用于把值转换成SQL语句中的字段值
var valueFunc = func(arg any) (any, error) {

	if arg == nil || reflect.ValueOf(arg).IsNil() {
		return "null", nil
	}
	switch arg := arg.(type) {
	case *string:
		return "'" + strings.ReplaceAll(*arg, "'", "\\'") + "'", nil
	case string:
		return "'" + strings.ReplaceAll(arg, "'", "\\'") + "'", nil
	case types.NullMarshalInterface:
		v, err := arg.MarshalJSON()
		return string(v[:]), err
	default:
		return arg, nil
	}

}

//参数必须是interface类型，基本类型会抛出错误
var rawValueFunc = func(arg any) (any, error) {

	if arg == nil || reflect.ValueOf(arg).IsNil() {
		return "null", nil
	}
	switch arg := arg.(type) {
	case *string:
		return strings.ReplaceAll(*arg, "'", "\\'"), nil
	case string:
		return strings.ReplaceAll(arg, "'", "\\'"), nil
	case types.NullMarshalInterface:
		v, err := arg.MarshalJSON()
		return string(v[:]), err
	default:
		return arg, nil
	}

}

//参数必须是interface类型，基本类型会抛出错误
var joinFunc = func(arg any) (any, error) {
	val := reflect.ValueOf(arg)
	if arg == nil || val.IsNil() {
		return "", nil
	}

	if val.Kind() != reflect.Slice {
		panic("参数不是slice类型")
	}
	len := val.Len()
	if len <= 0 {
		return "", nil
	}
	sb := &strings.Builder{}
	for i := 0; i < len; i++ {
		sb.WriteString(",")
		v := val.Index(i)
		// s,_:=rawValueFunc(v)
		sb.WriteString(fmt.Sprint(v))
	}
	s := sb.String()
	return s[1:], nil

}

func InitSqlTemplate() {
	tmpl, err := template.New("sqlTemplate").Funcs(
		template.FuncMap{
			"val":    valueFunc,
			"rawVal": rawValueFunc,
			"join":   joinFunc,
		}).ParseGlob("./resources/sql/*.sql.tmpl")
	if err != nil {
		panic(fmt.Sprintf("解释SQL模板文件错误！%v", err))
	}
	sqlTemplate = tmpl

}

//根据模板名称获取模板，并根据输入的数据执行模板
//TODO：缓存模板
func GetSqlString(t_name string, data any) string {
	tmpl := sqlTemplate.Lookup(t_name)
	if tmpl == nil {
		panic("找不到对应的sql模板")
	}
	//TODO:日志记录 info debug
	sb := &strings.Builder{}
	err := tmpl.Execute(sb, data)
	if err != nil {
		panic(err)
	}
	return sb.String()
}
