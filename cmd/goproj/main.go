//
// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.
//

package main

import (
  "errors"
  "flag"
  "fmt"
  "log"
  "os"
  "path/filepath"
  "strings"

  "github.com/demdxx/goproj"
)

const version = "2.1.2beta"
const author = "Dmitry Ponomarev <demdxx@gmail.com>"
const year = "2014-2015"

var (
  flagObserver bool = false
)

var help = map[string]string{
  "init":    "create project structure. goproj init [app-repository] <app-name>",
  "deps":    "list of dependencies",
  "list":    "list packages",
  "build":   "compile packages and dependencies",
  "clean":   "remove object files",
  "doc":     "run godoc on package sources",
  "fix":     "run go tool fix on packages",
  "fmt":     "run gofmt on package sources",
  "get":     "download and install packages and dependencies",
  "install": "compile and install packages and dependencies",
  "run":     "compile and run Go program",
  "test":    "test packages",
  "tool":    "run specified go tool",
  "vet":     "run go tool vet on packages",
  "version": "show goproj version",
  "info":    "print enviroment info",
  "go":      "return path for go",
  "help":    "show help or help [command]",
}

func init() {
  for k, v := range help {
    flag.Set(k, v)
  }
  flag.BoolVar(&flagObserver, "observer", false, "monitor and automatically restart")
  flag.Usage = printHelp
}

func main() {
  if len(os.Args) < 2 {
    printHelp()
    return
  }

  arg1 := os.Args[0]
  cmd := os.Args[1]
  os.Args = os.Args[1:]
  os.Args[0] = arg1

  if len(os.Args) > 1 && "-observer" == os.Args[1] {
    arg := os.Args[0]
    os.Args = os.Args[1:]
    os.Args[0] = arg
  }

  flag.Parse()
  pwd, err := os.Getwd()
  if nil != err {
    pwd = filepath.Dir(arg1)
  }

  args := os.Args[1:]

  switch cmd {
  case "init":
    cmdInitSolution(pwd, args)
    break
  case "deps":
    cmdDepList(pwd)
    break
  case "list":
    cmdProjList(pwd)
    break
  case "go":
    if len(args) > 0 {
      goproj.GoRun(args...)
    } else {
      fmt.Print(goproj.GoPath())
    }
    break
  case "path":
    printProjectPath(pwd)
    break
  case "solutionpath":
    printSolutionPath(pwd)
    break
  case "info":
    printInfo(pwd)
    break
  case "version":
    fmt.Printf("Goproj %s %s %s\n", version, author, year)
    break
  case "help":
    printHelp()
    break
  default:
    cmdExec(pwd, cmd, args)
  }
}

///////////////////////////////////////////////////////////////////////////////
/// Shell Commands
///////////////////////////////////////////////////////////////////////////////

func cmdInitSolution(pwd string, args []string) {
  if len(args) < 1 {
    log.Panicln("Err: invalid `init` command fromat")
    return
  }

  var project *goproj.Project = nil

  // Get or create solution for dir
  sol, err := goproj.SolutionForDir(pwd)
  if nil != err {
    log.Fatal(err)
  }

  if sol.IsGlobal {
    log.Fatal(errors.New("You can`t init project for global solution"))
  }

  // Add project
  if len(args) > 0 {
    project = goproj.ProjectFromUrl(args...)
    sol.AddProject(project)
  }

  // Init solution file structure
  if nil != project {
    if err := sol.InitFileStruct(project.Name); nil != err {
      log.Fatal(err)
      return
    }
  } else if err := sol.InitFileStruct(); nil != err {
    log.Fatal(err)
    return
  }

  // Save config files
  if err := sol.SaveConfig(); nil != err {
    log.Fatal(err)
    return
  }
}

func cmdDepList(pwd string) {
  sol := solution(pwd)
  if nil == sol {
    return
  }

  if nil != sol.Projects && len(sol.Projects) > 0 {
    for _, p := range sol.Projects {
      if nil != p.Deps && len(p.Deps) > 0 {
        for _, d := range p.Deps {
          fmt.Println(d.Url)
        }
      }
    }
  }
}

func cmdProjList(pwd string) {
  sol := solution(pwd)
  if nil == sol {
    return
  }

  if nil != sol.Projects && len(sol.Projects) > 0 {
    for _, p := range sol.Projects {
      if len(p.Url) > 0 {
        fmt.Println(p.Path, " => ", p.Url)
      } else {
        fmt.Println(p.Path)
      }
    }
  }
}

/**
 * Run project commands
 *
 * @param pwd
 * @param cmd
 */
func cmdExec(pwd, cmd string, args []string) {
  sol := solution(pwd)
  if nil == sol {
    return
  }

  var flags map[string]interface{} = nil // TODO parse all flags

  if err := sol.CmdExec(cmd, args, flags, flagObserver); nil != err {
    fmt.Println(err)
  }
}

func printHelp() {
  header := fmt.Sprintf("Goproj %s %s %s", version, author, year)
  fmt.Printf("%s\n%s\n", header, strings.Repeat("=", len(header)))

  for k, v := range help {
    fmt.Printf("% 10s - %s\n", k, v)
  }

  fmt.Print("\n")
  flag.PrintDefaults()
}

func printProjectPath(dir string) {
  path, _ := goproj.FindProjectDirFrom(dir)
  fmt.Print(path)
}

func printSolutionPath(dir string) {
  path, _ := goproj.FindSolutionDirFrom(dir)
  fmt.Print(path)
}

func printInfo(dir string) {
  env := goproj.SolutionEnv(dir)
  GOPATH := env["GOPATH"]

  for k, v := range goproj.SolutionEnv(dir) {
    fmt.Printf("%s=%s\n", k, v)
  }

  if len(GOPATH) < 1 {
    fmt.Println("Go environment for global projects not configured: GOPATH, GOBIN")
  }
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func solution(pwd string) *goproj.Solution {
  sol, err := goproj.SolutionFromDir(pwd)
  if nil == sol || nil != err {
    if nil != err {
      log.Print(err)
    }
    return nil
  }
  return sol
}
