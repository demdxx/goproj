// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package goproj

import (
  "log"
  "path/filepath"

  // "github.com/howeyc/fsnotify"
  "gopkg.in/fsnotify.v1"
)

func fsObserve(path string, change func() bool) (err error) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
    return
  }
  defer watcher.Close()

  done := make(chan bool)

  watcher.Add(path)
  files, _ := filepath.Glob(path + "/*")
  if nil != files && len(files) > 0 {
    for _, f := range files {
      watcher.Add(f)
    }
  } else {
    return
  }

  // Process events
  go func() {
    for {
      select {
      case ev := <-watcher.Events:
        if 0 != (ev.Op&fsnotify.Create) && ".go" == filepath.Ext(ev.Name) {
          if !change() {
            done <- true
            return
          }
        }
      case err = <-watcher.Errors:
        log.Println("error:", err)
      }
    }
  }()

  <-done
  return
}
