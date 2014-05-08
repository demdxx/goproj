package lib

import (
  "os/exec"
)

type goProc struct {
  Path string
}

func (g goProc) init() {
  g.Path = ""
}

func (g goProc) run(args ...string) error {
  cmd := exec.Command(g.Path, args...)
  return cmd.Run()
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
