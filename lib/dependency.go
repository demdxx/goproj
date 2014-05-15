// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "strings"
)

type Dependency struct {
  Path   string // Path at src dir
  Url    string // github.com/demdxx/goproj or git:https://github.com/demdxx/goproj#v2.0.1
  Config Config
}

///////////////////////////////////////////////////////////////////////////////
/// Commands
///////////////////////////////////////////////////////////////////////////////

// Get commands array
func (d *Dependency) Cmds() []interface{} {
  if cmds, ok := d.Config["cmd"]; ok {
    return cmds.([]interface{})
  }
  return nil
}

func (d *Dependency) Cmd(name, def string) string {
  if cmds := d.Cmds(); nil != cmds {
    if cmd, ok := cmds; ok {
      return cmd.(string)
    }
  }
  return def
}

// Shortcuts...

// @return {cmd} or ""
func (d *Dependency) CmdGet() string {
  _, cmd, url := PrepareCVSUrl(d.Url)
  cmd = strings.Replace(cmd, "{url}", url, 0)
  return d.Cmd("get", cmd)
}

// @return {cmd} or ""
func (d *Dependency) CmdBuild() string {
  return d.Cmd("build", "")
}

// @return {cmd} or ""
func (d *Dependency) CmdRun() string {
  return d.Cmd("run", "")
}
