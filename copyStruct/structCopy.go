package copyStruct

import (
  "reflect"
)

func StructCopy(structOne, structTwo interface{}) {

  if structOne == nil || structTwo == nil {
    return
  }

  ot := reflect.TypeOf(structOne)

  tt := reflect.TypeOf(structTwo)
  tv := reflect.ValueOf(structTwo)
  if ot.Kind() == reflect.Interface || ot.Kind() == reflect.Ptr {
    ot = ot.Elem()
  }
  if tt.Kind() == reflect.Interface || tt.Kind() == reflect.Ptr {
    tt = tt.Elem()
    tv = tv.Elem()
  }

  if ot.Kind() != reflect.Struct || tt.Kind() != reflect.Struct {
    return
  }

  structMap := make(map[string]reflect.Value)

  makeMapNameValue(structMap, structOne)

  for i := 0; i < tt.NumField(); i++ {
    s := tt.Field(i)
    if v, flag := structMap[s.Name]; flag {
      if s.Type == v.Type() {
        tv.Field(i).Set(v)
      }
    }
  }

}
func makeMapNameValue(structMap map[string]reflect.Value, interfaceStruct interface{}) {
  ot := reflect.TypeOf(interfaceStruct)
  ov := reflect.ValueOf(interfaceStruct)

  if ot.Kind() == reflect.Interface || ot.Kind() == reflect.Ptr {
    ot = ot.Elem()
    ov = ov.Elem()
  }

  var(
    v reflect.Value
    s reflect.StructField
  )

  for i := 0; i < ot.NumField(); i++ {
    s = ot.Field(i)
    v = ov.Field(i);
    if s.PkgPath == "" {
      if (v.Kind() == reflect.Struct) {
        makeMapNameValue(structMap, ov.Field(i).Interface())
      }else {
        if _, flag := structMap[s.Name]; flag {
          if isEmptyValue(v) {
            continue
          }
        }
        structMap[s.Name] = ov.Field(i)
      }
    }
  }
}

//func isEmptyValue(v reflect.Value) bool {
//  switch v.Kind() {
//  case reflect.String, reflect.Array:
//    return v.Len() == 0
//  case reflect.Map, reflect.Slice:
//    return v.Len() == 0 || v.IsNil()
//  case reflect.Bool:
//    return !v.Bool()
//  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//    return v.Int() == 0
//  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
//    return v.Uint() == 0
//  case reflect.Float32, reflect.Float64:
//    return v.Float() == 0
//  case reflect.Interface, reflect.Ptr:
//    return v.IsNil()
//  }
//  return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
//}