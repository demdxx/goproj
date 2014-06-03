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
  Owner  interface{} // Dependency high level owner
  Name   string      // Original config name
  Path   string      // Full Path at src dir
  Url    string      // github.com/demdxx/goproj or git:https://github.com/demdxx/goproj#v2.0.1
  Config Config
}

///////////////////////////////////////////////////////////////////////////////
/// Init
///////////////////////////////////////////////////////////////////////////////

// TODO Init enviroment before run any command
func (d *Dependency) UpdateEnv() {
  // TODO set custom environment from config
}

///////////////////////////////////////////////////////////////////////////////
/// Commands
///////////////////////////////////////////////////////////////////////////////

// Get commands array
func (d *Dependency) Cmds() map[string]interface{} {
  if cmds, ok := d.Config["cmd"]; ok {
    switch cmds.(type) {
    case map[interface{}]interface{}:
      return ToStringMap(cmds)
    case map[string]interface{}:
      return cmds.(map[string]interface{})
    }
  }
  return nil
}

func (d *Dependency) Cmd(name string, def interface{}) interface{} {
  if cmds := d.Cmds(); nil != cmds {
    if cmd, ok := cmds[name]; ok {
      return cmd
    }
  }
  return def
}

// Shortcuts...

// @return {cmd} or ""
func (d *Dependency) CmdGet() interface{} {
  _, cmd, url := PrepareCVSUrl(d.Url)
  if len(url) < 1 {
    return nil
  }
  cmd = strings.Replace(cmd, "{url}", url, -1)
  return d.Cmd("get", cmd)
}

// @return {cmd} or ""
func (d *Dependency) CmdBuild() interface{} {
  return d.Cmd("build", "")
}

// @return {cmd} or ""
func (d *Dependency) CmdRun() interface{} {
  return d.Cmd("run", "")
}

func (d *Dependency) CmdTest() interface{} {
  return d.Cmd("test", "{go} test -v {flags} {app}")
}
