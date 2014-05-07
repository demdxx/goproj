// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014

package main

import (
  "fmt"
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

func main() {
  printHelp()
}

func printHelp() {
  fmt.Printf("Goproj %s %s\n", version, author)
  fmt.Print("================================================\n")

  for k, v := range help {
    fmt.Printf("% 10s: %s\n", k, v)
  }
}
