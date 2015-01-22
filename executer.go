// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package goproj

import (
  "bytes"
  "errors"
  "fmt"
  "log"
  "os"
  "os/exec"
  "strconv"
  "strings"

  "github.com/demdxx/gocast"
)

type CommandExecutor interface {
  UpdateEnv()

  Cmds() map[string]interface{}
  Cmd(name string, def interface{}) interface{}

  // Shortcuts...
  // @return {cmd} or ""

  CmdGet() interface{}
  CmdBuild() interface{}
  CmdInstall() interface{}
  CmdRun() interface{}
  CmdTest() interface{}
}

///////////////////////////////////////////////////////////////////////////////
/// Prepare
///////////////////////////////////////////////////////////////////////////////

func getSolution(e interface{}) *Solution {
  switch e.(type) {
  case Dependency:
    return e.(Dependency).Solution()
    break
  case *Dependency:
    return e.(*Dependency).Solution()
    break
  case Project:
    return e.(Project).Solution()
    break
  case *Project:
    return e.(*Project).Solution()
    break
  }
  return nil
}

func getSolutionPath(e interface{}) string {
  switch e.(type) {
  case *Dependency:
    return e.(*Dependency).SolutionPath()
    break
  case *Project:
    return e.(*Project).SolutionPath()
    break
  }
  return ""
}

func getPath(e interface{}) string {
  switch e.(type) {
  case Dependency:
    return e.(Dependency).Path
    break
  case *Dependency:
    return e.(*Dependency).Path
    break
  case Project:
    return e.(Project).Path
    break
  case *Project:
    return e.(*Project).Path
    break
  }
  return ""
}

func getFullPath(e interface{}) string {
  sol := getSolution(e)
  return fmt.Sprintf("%s/src/%s", sol.Path, getPath(e))
}

func getApp(e interface{}) string {
  switch e.(type) {
  case Dependency:
    return e.(Dependency).Name
    break
  case *Dependency:
    return e.(*Dependency).Name
    break
  case Project:
    return e.(Project).Name
    break
  case *Project:
    return e.(*Project).Name
    break
  }
  return ""
}

func prapareFlags(flags map[string]interface{}) string {
  buf := bytes.NewBuffer(nil)
  if nil != flags {
    for k, v := range flags {
      buf.WriteByte('&')
      buf.WriteString(k)
      buf.WriteByte('=')
      switch v.(type) {
      case string:
        buf.WriteString(v.(string))
        break
      case int:
      case int32:
      case int64:
        buf.WriteString(strconv.Itoa(v.(int)))
        break
      default:
        buf.WriteString(fmt.Sprintf("%v", v))
      }
    }
  }
  return buf.String()
}

func prepareCommand(e CommandExecutor, cmd interface{}, flags map[string]interface{}) (interface{}, error) {
  switch cmd.(type) {
  case string:
    s := cmd.(string)
    params, err := gocast.ToStringMap(e.Cmds(), "")
    if nil != err {
      log.Panic(err)
      return "", nil
    } else if len(params) > 0 {
      for k, v := range params {
        s = strings.Replace(s, "{"+k+"}", v, -1)
      }
    }

    params = make(map[string]string)
    params["flags"] = prapareFlags(flags)
    params["solutionpath"] = getSolutionPath(e)
    params["fullpath"] = getFullPath(e)
    params["path"] = getPath(e)
    params["app"] = getApp(e)
    params["go"] = goproc.Path

    for k, v := range params {
      s = strings.Replace(s, "{"+k+"}", v, -1)
    }
    return s, nil
  }
  return "", errors.New("Prepare command failed")
}

///////////////////////////////////////////////////////////////////////////////
/// Exec
///////////////////////////////////////////////////////////////////////////////

func run(e CommandExecutor, command string) error {
  cmd, err := runCommand(e, command)
  if nil != err {
    return err
  }
  return cmd.Wait()
}

func runCommand(e CommandExecutor, command string) (*exec.Cmd, error) {
  e.UpdateEnv()
  fmt.Println(">", command)
  cmd := exec.Command("sh", "-c", command)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  if err := cmd.Start(); nil != err {
    return nil, err
  }
  return cmd, nil
}

func runObserver(e CommandExecutor, command, path string) error {
  cmd, err := runCommand(e, command)
  if nil == err {
    done := make(chan error, 1)
    go func() { done <- cmd.Wait() }()

    fmt.Println("fsObserve")
    fsObserve(path, func() bool { // Restart command here
      fmt.Println("Restart...")
      <-done // Exit goroutine
      if err := cmd.Process.Kill(); err != nil {
        log.Println("Failed to kill: ", err)
      }

      // Restart command
      done = make(chan error, 1)
      cmd, err := runCommand(e, command)
      if nil == err {
        go func() { done <- cmd.Wait() }()
      } else {
        log.Println("Failed run command: ", err)
      }
      return true
    })
  }
  return err
}

///////////////////////////////////////////////////////////////////////////////
/// Actions
///////////////////////////////////////////////////////////////////////////////

func execute(e CommandExecutor, cmd string, flags map[string]interface{}, observe bool) error {
  command := getCmd(e, cmd)

  if !isEmpty(command) {
    // Execute command
    switch command.(type) {
    case string:
      // Prepare command
      var err error
      if command, err = prepareCommand(e, command, flags); nil != err {
        return err
      }
      if observe {
        return runObserver(e, command.(string), getSolutionPath(e))
      }
      return run(e, command.(string))
    case []interface{}:
      var err error
      var cmd interface{}
      for _, c := range command.([]interface{}) {
        // Prepare command
        if cmd, err = prepareCommand(e, c, flags); nil != err {
          return err
        }

        if observe {
          if err = runObserver(e, cmd.(string), getSolutionPath(e)); nil != err {
            return err
          }
        } else if err = run(e, cmd.(string)); nil != err {
          return err
        }
      }
      break
    case []string:
      var err error
      var cmd interface{}
      for _, c := range command.([]string) {
        // Prepare command
        if cmd, err = prepareCommand(e, c, flags); nil != err {
          return err
        }

        if observe {
          if err = runObserver(e, cmd.(string), getSolutionPath(e)); nil != err {
            return err
          }
        } else if err = run(e, cmd.(string)); nil != err {
          return err
        }
      }
      break
    }
  }

  return nil // errors.New(fmt.Sprintf("Unsupport command: %s", cmd))
}

func getCmd(e CommandExecutor, cmd string) interface{} {
  switch cmd {
  case "get":
    return e.CmdGet()
  case "build":
    return e.CmdBuild()
  case "install":
    return e.CmdInstall()
  case "run":
    return e.CmdRun()
  case "test":
    return e.CmdTest()
  default:
    return e.Cmd(cmd, "")
  }
}