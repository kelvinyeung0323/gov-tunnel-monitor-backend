package types

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type NullMarshalInterface interface {
	MarshalJSON() ([]byte, error)
}

//因为json 对time.Time解释有bug,这里自定义Time类型序列化和反序列化处理
type Time time.Time

func (v Time) MarshalJSON() ([]byte, error) {
	formatted := "\"" + time.Time(v).Format(TimeFormat) + "\""
	return []byte(formatted), nil

}

func (v *Time) UnmarshalJSON(data []byte) error {
	// TODO:提出来参数配置
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(TimeFormat, strings.ReplaceAll(string(data), "\"", ""), loc)
	if err != nil {
		return err
	}
	b := Time(t)
	*v = b
	return nil
}

type NullTime sql.NullTime

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}

	return nil
}

type NullString sql.NullString

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

type NullInt16 sql.NullInt16

func (v NullInt16) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int16)
	} else {
		return json.Marshal(nil)
	}
}

type NullInt32 sql.NullInt32

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

type NullInt64 sql.NullInt64

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

type NullBool sql.NullBool

func (v NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(nil)
	}
}

type NullByte sql.NullByte

func (v NullByte) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Byte)
	} else {
		return json.Marshal(nil)
	}
}

type NullFloat64 sql.NullFloat64

func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	} else {
		return json.Marshal(nil)
	}
}
