// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "fmt"
  "os/exec"
)

type goProc struct {
  Path string
}

func (g *goProc) init() {
  g.Path = "go" // TODO detect `go` path
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

func GoPath() string {
  return goproc.Path
}

func GoRun(args ...string) error {
  return goproc.run(args...)
}
