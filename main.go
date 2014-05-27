// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package main

import (
  "flag"
  "fmt"
  "goproj/lib"
  "log"
  "os"
  "path/filepath"
  "strings"
)

const version = "2.0.0alpha"
const author = "Dmitry Ponomarev <demdxx@gmail.com>"
const year = "2014"

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
  flag.Usage = printHelp
}

func main() {
  flag.Parse()
  pwd, err := os.Getwd()
  if nil != err {
    pwd = filepath.Dir(os.Args[0])
  }

  if flag.NArg() > 0 {
    switch os.Args[1] {
    case "init":
      cmdInitSolution(pwd, os.Args[2:])
      break
    case "deps":
      cmdDepList(pwd)
      break
    case "list":
      cmdProjList(pwd)
      break
    case "go":
      args := os.Args[2:]
      if len(args) > 0 {
        lib.GoRun(os.Args[2:]...)
      } else {
        fmt.Print(lib.GoPath())
      }
      break
    case "path":
      printProjectPath(pwd)
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
      cmdExec(pwd, os.Args[1], flag.Args()[1:])
    }
  }
}

///////////////////////////////////////////////////////////////////////////////
/// Commands
///////////////////////////////////////////////////////////////////////////////

func cmdInitSolution(pwd string, args []string) {
  if len(args) < 1 {
    log.Panicln("Err: invalid `init` command fromat")
    return
  }

  sol, _ := lib.SolutionFromDir(pwd)
  if nil == sol {
    sol = &lib.Solution{}
    sol.Init(pwd)
  }

  // Add project
  if len(args) > 0 {
    sol.AddProject(lib.ProjectFromUrl(args...))
  }

  // Init file structure
  if err := sol.InitFileStruct(); nil != err {
    log.Print(err)
    return
  }

  // Save configs
  if err := sol.SaveConfig(); nil != err {
    log.Print(err)
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
  if err := sol.CmdExec(cmd, args, flags); nil != err {
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
  projdir, _ := lib.FindProjectDirFrom(dir)
  fmt.Print(projdir)
}

func printInfo(dir string) {
  for k, v := range lib.SolutionEnv(dir) {
    fmt.Printf("%s=%s\n", k, v)
  }
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func solution(pwd string) *lib.Solution {
  sol, err := lib.SolutionFromDir(pwd)
  if nil == sol || nil != err {
    if nil != err {
      log.Print(err)
    }
    return nil
  }
  return sol
}
