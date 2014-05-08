package lib

import (
  "fmt"
  "os/exec"
)

type goProc struct {
  Path string
}

func (g *goProc) init() {
  g.Path = "go"
}

func (g goProc) run(args ...string) (err error) {
  out, err := exec.Command(g.Path, args...).Output()
  if err != nil {
    return
  }
  fmt.Printf("%s", out)
  return
}

///////////////////////////////////////////////////////////////////////////////
/// go stand alone
///////////////////////////////////////////////////////////////////////////////

var goproc goProc

func init() {
  goproc.init()
}

func GoRun(args ...string) error {
  return goproc.run(args...)
}
