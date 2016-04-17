package vaild

import (
	"reflect"
//"fmt"
//"regexp"
	"strings"
)

var Errors = map[string]string{
	"required":"不能为空",
}

type ErrorVaildStruct struct {
	ErrorName string
	VaildName string
	Msg string
}

func NewVaildError(errorName,vaildName,msg string) error {
	return ErrorVaildStruct{
		ErrorName:errorName,
		VaildName:vaildName,
		Msg:msg,
	}
}

func (this ErrorVaildStruct) Error() string {
	return this.Msg
}



func ValidateStruct(su interface{}) error {
	if su == nil {
		return nil
	}

	tof := reflect.TypeOf(su)
	vof := reflect.ValueOf(su)

	if tof.Kind() == reflect.Interface || tof.Kind() == reflect.Ptr {
		tof = tof.Elem()
		vof = vof.Elem()
	}
	// we only accept structs
	if tof.Kind() != reflect.Struct {
		return nil
	}

	var t reflect.StructField
	var p reflect.Value
	var tagStr string
	for i := 0; i <= tof.NumField(); i++ {

		t = tof.Field(i)
		tagStr = t.Tag.Get("vaild") //获取 tag

		if tagStr != "" {
			p = vof.Field(i)
			err := vaildSwitch(t.Name,p,tagStr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateStructMap(su interface{}, cond map[string]string) error {
	if len(cond) == 0 || su == nil {
		return nil
	}

	tof := reflect.TypeOf(su)
	vof := reflect.ValueOf(su)

	if tof.Kind() == reflect.Interface || tof.Kind() == reflect.Ptr {
		tof = tof.Elem()
		vof = vof.Elem()
	}
	// we only accept structs
	if tof.Kind() != reflect.Struct {
		return nil
	}

	var t reflect.StructField
	var p reflect.Value
	for i := 0; i <= tof.NumField(); i++ {

		t = tof.Field(i)

		if v, flag := cond[t.Name]; flag {
			if v != "" {
				p = vof.Field(i)

				err := vaildSwitch(t.Name, p, v)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func vaildSwitch(name string, p reflect.Value, cond string) error {


	conds := strings.Split(cond," ")

	for _, v := range conds {

		switch v {
		//是否为空
		case "required":
			if isEmptyValue(p) {
				return NewVaildError(name,"required",Errors["required"])
			}

		}
	}

	return nil
}

//func isMobile(v reflect.Value) bool{
//	if v.Kind != reflect.String {
//		return false
//	}
//
//
//
//
//}

//是否为空
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}