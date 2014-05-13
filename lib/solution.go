// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "errors"
  "fmt"
)

type Solution struct {
  Path     string
  Projects []*Project
  Config   Config
}

func SolutionFromFile(filepath string) (sol *Solution, err error) {
  if err = sol.Config.InitFromFile(filepath); nil == err {
    err = sol.Init()
  }
  return
}

func (sol *Solution) Init() error {
  if nil == sol.Config || len(sol.Config) < 1 {
    return errors.New("Project not inited")
  }

  if nil != sol.Projects {
    sol.Projects = nil
  }

  // Each config
  for dir, conf := range sol.Config {
    proj, err := ProjectFromFile(dir, conf.(Config))
    if nil == err {
      err = sol.AddProject(proj)
    }
    if nil != err {
      return err
    }
  }
  return nil
}

func (sol *Solution) AddProject(p *Project) error {
  if nil == sol.Projects {
    sol.Projects = make([]*Project, 0)
  }
  sol.Projects = append(sol.Projects, p)
  return nil
}

// Init FS struct
//
// bin/
// pkg/
// src/
// .gosolution
func (sol *Solution) InitFileStruct() error {
  if len(sol.Path) < 1 {
    return errors.New("Solution path not defined")
  }
  if err := makeDir(sol.Path); nil != err {
    return err
  }

  // Create dirs
  if err := makeDir(fmt.Sprintf("%s/bin", sol.Path)); nil != err {
    return err
  }
  path = fmt.Sprintf("%s/bin", sol.Path)
  if err := makeDir(fmt.Sprintf("%s/pkg", sol.Path)); nil != err {
    return err
  }
  if err := makeDir(fmt.Sprintf("%s/src", sol.Path)); nil != err {
    return err
  }

  // Create solution
  sol.SaveConfig()
}

///////////////////////////////////////////////////////////////////////////////
/// Other
///////////////////////////////////////////////////////////////////////////////

func FindSolutionFrom(dir string) (string, error) {
  return findParentDirWithFile(sir, ".gosolution")
}

func HasSolution(dir string) bool {
  return isFile(fmt.Sprintf("%s/.gosolution", dir))
}
