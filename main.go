//
// @project goproj
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package main

import (
  "fmt"
  "os"

  "gopkg.in/alecthomas/kingpin.v2"

  "./cmd"
  "./modules"
)

var (
  app              = kingpin.New("goproj", "Goproj development tool")
  cmdGo            = app.Command("go", "Run raw go command")
  cmdInit          = app.Command("init", "Init project")
  flagInitGlobal   = cmdInit.Flag("global", "Init as singletone global project").Bool()
  argInitUrl       = cmdInit.Arg("url", "Repository URL").String()
  argInitProject   = cmdInit.Arg("project", "Project name").String()
  cmdDeps          = app.Command("dependecy", "Control depenencies").Alias("dep")
  cmdDepList       = cmdDeps.Command("list", "Depenency list").Default()
  cmdDepAdd        = cmdDeps.Command("add", "Add depenency by URL")
  flagDepAddTarget = cmdDepAdd.Flag("target", "Config target section").String()
  argDepAddDep     = cmdDepAdd.Arg("dependency", "Dependency URL").Required().String()
  argDepAddProject = cmdDepAdd.Arg("project", "Dependency URL").String()
  cmdDepGet        = cmdDeps.Command("get", "Download all depenencies")
  flagDepGetTarget = cmdDepGet.Flag("target", "Config target section").String()
  argDepGetProject = cmdDepGet.Arg("project", "Project name").String()
  cmdEnv           = app.Command("env", "Print project environment")
)

func main() {
  command := kingpin.MustParse(app.Parse(os.Args[1:]))
  sol := modules.GetSolution()

  switch command {
  case cmdEnv.FullCommand():
    cmd.EnvPrint(sol)
    break
  case cmdGo.FullCommand():
    fmt.Println(os.Args[2:])
    break
  case cmdDepList.FullCommand():
    fmt.Println(cmdDepList.FullCommand(), os.Args[2:])
    break
  }
}
