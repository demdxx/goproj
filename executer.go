//
// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.
//

package goproj

import (
  "bytes"
  "errors"
  "fmt"
  "log"
  "os"
  "os/exec"
  "os/signal"
  "path/filepath"
  "runtime/pprof"
  "strconv"
  "strings"
  "syscall"

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
  path := getPath(e)
  if filepath.IsAbs(path) {
    return path
  }
  sol := getSolution(e)
  return fmt.Sprintf("%s/src/%s", sol.Path, path)
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

func prapareArgs(args []string) string {
  if nil == args {
    return ""
  }
  return strings.Join(args, " ")
}

func prepareCommand(e CommandExecutor, cmd interface{}, args []string, flags map[string]interface{}) (interface{}, error) {
  switch cmd.(type) {
  case string:
    s := cmd.(string)
    params, err := gocast.ToStringMap(e.Cmds(), "", false)
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
    params["args"] = prapareArgs(args)
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
  if nil == err {
    // capture ctrl+c and stop CPU profiler
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func() {
      for _ = range c {
        killCmd(cmd)
        pprof.StopCPUProfile()
        os.Exit(1)
      }
    }()
    return cmd.Wait()
  }
  return err
}

func runCommand(e CommandExecutor, command string) (*exec.Cmd, error) {
  e.UpdateEnv()
  fmt.Println(">", "\033[0;32m"+command+"\033[0m")
  cmd := exec.Command("sh", "-c", command)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  if err := cmd.Start(); nil != err {
    return nil, err
  }
  return cmd, nil
}

func killCmd(cmd *exec.Cmd) {
  fmt.Println("Kill Proc", cmd)

  if nil != cmd {
    var err error
    pid := cmd.Process.Pid
    if err = syscall.Kill(pid, syscall.SIGTERM); nil != err {
      log.Println("Failed to kill:", err)
    }

    gpid, _ := syscall.Getpgid(pid)
    if err := syscall.Kill(-gpid, 15); nil != err {
      log.Println("Failed to kill process group:", gpid, err)
    }

    if err = cmd.Process.Signal(os.Kill); nil != err {
      log.Println("Failed to kill:", err)
    } else if _, err = cmd.Process.Wait(); nil != err {
      log.Println("Failed process wait:", err)
    }
    cmd = nil
  }
}

func runObserver(e CommandExecutor, command, path string) error {
  cmd, err := runCommand(e, command)
  if nil == err {
    // capture ctrl+c and stop CPU profiler
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func() {
      for _ = range c {
        killCmd(cmd)
        pprof.StopCPUProfile()
        os.Exit(1)
      }
    }()

    // Begin observer for folders in path
    fsObserve(path, func() bool {
      fmt.Println("Restart...")
      killCmd(cmd)

      // Restart command
      cmd, err = runCommand(e, command)
      if nil != err {
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

func execute(e CommandExecutor, cmd string, args []string, flags map[string]interface{}, observe bool) error {
  command := getCmd(e, cmd)
  path := getFullPath(e)

  if !isEmpty(command) {
    // Execute command
    switch command.(type) {
    case string:
      // Prepare command
      var err error
      if command, err = prepareCommand(e, command, args, flags); nil != err {
        return err
      }
      if observe {
        return runObserver(e, command.(string), path)
      }
      return run(e, command.(string))
    case []interface{}:
      var err error
      var cmd interface{}
      for _, c := range command.([]interface{}) {
        // Prepare command
        if cmd, err = prepareCommand(e, c, args, flags); nil != err {
          return err
        }

        if observe {
          if err = runObserver(e, cmd.(string), path); nil != err {
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
        if cmd, err = prepareCommand(e, c, args, flags); nil != err {
          return err
        }

        if observe {
          if err = runObserver(e, cmd.(string), path); nil != err {
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

// Get command
//
// @param executer
// @param cmd
// @return Prepared Command
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
