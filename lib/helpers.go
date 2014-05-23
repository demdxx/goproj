// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "reflect"
  "strings"
)

// From src/pkg/encoding/json.
func isEmptyValue(v reflect.Value) bool {
  switch v.Kind() {
  case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
    return v.Len() == 0
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
  return false
}

func isEmpty(v interface{}) bool {
  return nil == v || isEmptyValue(reflect.ValueOf(v))
}

func indexOfStringSlice(slice []string, s string) int {
  if nil != slice {
    for i, v := range slice {
      if s == v {
        return i
      }
    }
  }
  return -1
}

var urlSufixes = []string{
  ".git", ".hg", ".svn", ".bzr",
}

func IsUrl(url string) bool {
  arr := strings.Split(url, "://")
  if len(arr) < 2 {
    for _, fx := range urlSufixes {
      if strings.HasSuffix(url, fx) {
        return true
      }
    }
    return false
  }
  return true
}
