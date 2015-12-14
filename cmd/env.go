package cmd

import (
  "fmt"

  "../modules"
)

func EnvPrint(sol *modules.Solution) error {
  if nil == sol {
    fmt.Println("Project not found")
  }
  return nil
}
