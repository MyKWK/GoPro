package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// 将data映射到obj之中
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ { // 遍历结构体所有字段
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		// 如果字段标签是 sql:"id" ，则获取data["id"]的值 ，赋值给value
		name := objValue.Type().Field(i).Name       // 获取字段名称
		structFieldType := objValue.Field(i).Type() // 获取该字段类型，比如ID，ProductName
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			val, err = TypeConversion(value, structFieldType.Name())
			// 这里返回的是Value类型的变量
			if err != nil {
			}
		}
		objValue.FieldByName(name).Set(val) // 为对象赋值
	}
}

// 将string类型数值转换为对应类型
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}
	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
