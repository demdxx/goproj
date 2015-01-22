package goproj

import (
  "errors"
)

var (
  ErrorGoEnvironmentNotConfigured = errors.New("For global solution project need set enviroment GOPATH")
  ErrorGlobalSolutionSaveInvalid  = errors.New("You can`t save solution file for global type")
  ErrorSolutionPath               = errors.New("Invalid solution path")
)
