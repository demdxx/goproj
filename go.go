//
// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.
//

package goproj

import (
  "fmt"
  "os"
  "os/exec"
)

type goProc struct {
  Path string
}

func (g *goProc) init() {
  g.Path = "go"
}

func (g goProc) run(args ...string) error {
  fmt.Println(args)
  cmd := exec.Command(g.Path, args...)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  if err := cmd.Start(); nil != err {
    return err
  }
  return cmd.Wait()
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
