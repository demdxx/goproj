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

const version = "2.0.0"
const author = "Dmitry Ponomarev <demdxx@gmail.com>"

var help = map[string]string{
  "init":    "create project structure. goptoj init [folder] <name>",
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
  fmt.Println(flag.Args())
  pwd, err := os.Getwd()
  if nil != err {
    pwd = filepath.Dir(os.Args[0])
  }

  if flag.NArg() > 0 {
    switch os.Args[1] {
    case "init":
      cmdInitSolution(pwd)
      break
    case "deps":
      cmdDepList(pwd)
      break
    case "list":
      cmdProjList(pwd)
      break
    case "go":
      lib.GoRun(os.Args[2:]...)
      break
    case "info":
      printInfo(pwd)
      break
    }
  }
}

///////////////////////////////////////////////////////////////////////////////
/// Commands
///////////////////////////////////////////////////////////////////////////////

func cmdInitSolution(pwd string) {
  sol, err := lib.SolutionFromDir(pwd)
  if nil == sol || nil != err {
    if nil != err {
      log.Print(err)
    }
    return
  }
}

func cmdDepList(pwd string) {
  sol, err := lib.SolutionFromDir(pwd)
  if nil == sol || nil != err {
    if nil != err {
      log.Print(err)
    }
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
  sol, err := lib.SolutionFromDir(pwd)
  if nil == sol || nil != err {
    if nil != err {
      log.Print(err)
    }
    return
  }

  if nil != sol.Projects && len(sol.Projects) > 0 {
    for _, p := range sol.Projects {
      fmt.Println(p.Path, " => ", p.Url)
    }
  }
}

func printHelp() {
  header := fmt.Sprintf("Goproj %s %s", version, author)
  fmt.Printf("%s\n%s\n", header, strings.Repeat("=", len(header)))

  for k, v := range help {
    fmt.Printf("% 10s - %s\n", k, v)
  }

  fmt.Print("\n")
  flag.PrintDefaults()
}

func printInfo(dir string) {
  for k, v := range lib.SolutionEnv(dir) {
    fmt.Println(k, " = ", v)
  }
}
