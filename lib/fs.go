package lib

import (
  "errors"
  "fmt"
  "os"
  "path"
  "path/filepath"
)

func findParentDirWithFile(dir, name string) (string, error) {
  if dir, err := filepath.Abs(dir); nil != err {
    return "", err
  }
  for {
    filename := dir + "/" + name
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
