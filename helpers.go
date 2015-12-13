// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package goproj

import (
  "errors"
  "io"
  "net/url"
  "os"
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

func IsUrl(_url string) bool {
  arr := strings.Split(_url, "://")
  if len(arr) < 2 {
    for _, fx := range urlSufixes {
      if strings.HasSuffix(_url, fx) {
        return true
      }
    }
    return false
  }
  return true
}

func PathFromUrl(_url string) (string, error) {
  if !IsUrl(_url) {
    return "", errors.New("Invalid url format")
  }
  u, err := url.Parse(_url)
  if nil != err {
    return "", err
  }
  for _, fx := range urlSufixes {
    if strings.HasSuffix(u.Path, fx) {
      return u.Host + u.Path[:len(u.Path)-len(fx)], nil
    }
  }
  return u.Host + u.Path, nil
}

// func ToStringMap(data interface{}) map[string]interface{} {
//   switch data.(type) {
//   case map[string]interface{}:
//     return data.(map[string]interface{})
//   case map[interface{}]interface{}:
//     m := make(map[string]interface{})
//     for k, v := range data.(map[interface{}]interface{}) {
//       switch k.(type) {
//       case string:
//         switch v.(type) {
//         case map[interface{}]interface{}:
//           m[k.(string)] = ToStringMap(v)
//           break
//         default:
//           m[k.(string)] = v
//         }
//         break
//       }
//     }
//     return m
//   }
//   return nil
// }

// exists returns whether the given file or directory exists or not
//
// @param path
// @return bool
func isFsExists(path string) bool {
  _, err := os.Stat(path)
  if err == nil {
    return true
  }
  return !os.IsNotExist(err)
}

// Check is directory is empty
func isDirEmpty(dirname string) bool {
  f, err := os.Open(dirname)
  if err != nil {
    return false
  }
  defer f.Close()

  _, err = f.Readdir(1)
  return err == io.EOF
}
