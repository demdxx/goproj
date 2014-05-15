// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "errors"
  "fmt"
  "os"
  "path/filepath"
)

func findParentDirWithFile(dir, name string) (string, error) {
  {
    var err error
    if dir, err = filepath.Abs(dir); nil != err {
      return "", err
    }
  }
  for {
    filename := dir + "/" + name
    fmt.Println("findParentDirWithFile", filename)
    if info, err := os.Stat(filename); err == nil && !info.IsDir() {
      return dir, nil
    }

    // Get parent dir
    {
      ndir := filepath.Dir(dir)
      if ndir == dir || len(dir) < 2 {
        break
      }
      dir = ndir
    }
  }
  return "", errors.New(fmt.Sprintf("%s doesn't exists", name))
}

func isDir(fullpath string) (bool, error) {
  info, err := os.Stat(fullpath)
  if err == nil && info.IsDir() {
    return true, nil
  }
  return false, err
}

func isFile(fullpath string) (bool, error) {
  res, err := isDir(fullpath)
  if nil == err {
    return !res, err
  }
  return res, err
}

func makeDir(fullpath string) error {
  if info, err := os.Stat(fullpath); nil == info && err != nil {
    if err = os.MkdirAll(fullpath, 0755); nil != err {
      return err
    }
  }
  return nil
}
