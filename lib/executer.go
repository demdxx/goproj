package lib

import (
  "bytes"
  "errors"
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "strings"
)

type CommandExecutor interface {
  UpdateEnv()

  Cmds() map[string]interface{}
  Cmd(name string, def interface{}) interface{}

  // Shortcuts...
  // @return {cmd} or ""

  CmdGet() interface{}
  CmdBuild() interface{}
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
  // fmt.Println("getFullPath", sol.Path, getPath(e))
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
    var s string
    s = strings.Replace(cmd.(string), "{flags}", prapareFlags(flags), -1)
    s = strings.Replace(s, "{solutionpath}", getSolutionPath(e), -1)
    s = strings.Replace(s, "{fullpath}", getFullPath(e), -1)
    s = strings.Replace(s, "{path}", getPath(e), -1)
    s = strings.Replace(s, "{app}", getApp(e), -1)
    s = strings.Replace(s, "{go}", goproc.Path, -1)
    return s, nil
    break
  }
  return "", errors.New("Prepare command failed")
}

func run(e CommandExecutor, command string) error {
  e.UpdateEnv()
  fmt.Println(">", command)
  cmd := exec.Command("sh", "-c", command)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  if err := cmd.Start(); nil != err {
    return err
  }
  return cmd.Wait()
}

///////////////////////////////////////////////////////////////////////////////
/// Actions
///////////////////////////////////////////////////////////////////////////////

func execute(e CommandExecutor, cmd string, flags map[string]interface{}) error {
  var command interface{} = nil
  switch cmd {
  case "get":
    command = e.CmdGet()
    break
  case "build":
    command = e.CmdBuild()
    break
  case "run":
    command = e.CmdRun()
    break
  case "test":
    command = e.CmdTest()
    break
  default:
    command = e.Cmd(cmd, "")
    break
  }

  if !isEmpty(command) {
    // Execute command
    switch command.(type) {
    case string:
      // Prepare command
      var err error
      if command, err = prepareCommand(e, command, flags); nil != err {
        return err
      }
      return run(e, command.(string))
      break
    case []interface{}:
      var err error
      var cmd interface{}
      for _, c := range command.([]interface{}) {
        // Prepare command
        if cmd, err = prepareCommand(e, c, flags); nil != err {
          return err
        }
        if err = run(e, cmd.(string)); nil != err {
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
        if err = run(e, cmd.(string)); nil != err {
          return err
        }
      }
      break
    }
  }

  return nil // errors.New(fmt.Sprintf("Unsupport command: %s", cmd))
}
