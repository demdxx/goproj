// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "errors"
  "fmt"
  "gopkg.in/v1/yaml"
  "io/ioutil"
  "os"
  "path/filepath"
)

type Solution struct {
  IsGlobal bool // We don't have local solution
  Path     string
  Projects []*Project
  Config   Config
}

func SolutionFromDir(dir string) (sol *Solution, err error) {
  ndir, err := FindSolutionDirFrom(dir)
  if nil != err {
    // We don't have .gosolution then check .goproj
    pdir, err := findParentDirWithFile(sir, ".goproj")
    if nil != err {
      return nil, err
    }

    // Init global solution
    if GOPATH, ok := os.Environ()["GOPATH"]; ok {
      err = nil
      sol := new(Solution)
      sol.IsGlobal = true
      sol.Config["projects"] = map[string]interface{}{
        pdir: map[string]interface{}{}, // Empty config
      }
      return sol, sol.Init(GOPATH)
    }
  }

  // Init solution from file
  sol, err = SolutionFromFile(fmt.Sprintf("%s/.gosolution", ndir))
  return
}

func SolutionFromFile(fpath string) (sol *Solution, err error) {
  var conf Config
  if err = conf.InitFromFile(fpath); nil == err {
    sol := new(Solution)
    sol.Config = conf
    err = sol.Init(filepath.Dir(fpath))
  }
  return
}

func (sol *Solution) Init(path string) error {
  if len(path) > 0 {
    sol.Path, err = filepath.Abs(path)
    if nil != err {
      return err
    }
  }

  if nil == sol.Config || len(sol.Config) < 1 {
    return errors.New("Project not inited")
  }

  if nil != sol.Projects {
    sol.Projects = nil
  }

  // Each config
  if projects, ok := sol.Config["projects"]; ok {
    if nil != projects {
      switch projects.(type) {
      case map[string]interface{}:
        for dir, conf := range projects {
          proj, err := ProjectFromFile(dir, conf.(Config))
          if nil == err {
            err = sol.AddProject(proj)
          }
          if nil != err {
            return err
          }
        }
      default:
        return errors.New("Config has invalid format in conf.projects section")
      }
    }
  }
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
  return sol.SaveConfig()
}

// Save soluton project
//
// @return nil or error
func (sol *Solution) SaveConfig() error {
  projects := make(map[string]interface{})

  for _, p := range sol.Projects {
    projects[p.Path] = map[string]interface{}{
      "url": p.Url,
    }
  }

  // Murshal config to YAML
  data, err := yaml.Marshal(map[string]interface{}{
    "projects": projects,
  })
  if nil != err {
    return err
  }

  // Store file
  return ioutil.WriteFile(fmt.Sprintf("%s/.gosolution", dir), data, os.FileMode)
}

///////////////////////////////////////////////////////////////////////////////
/// Projects
///////////////////////////////////////////////////////////////////////////////

// Init solution projects by repos or crreate project structure
//
// @return nil or error
func (sol *Solution) InitProjects() error {
  for _, p := range sol.Projects {
    if err := p.Init(); nil != err {
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

///////////////////////////////////////////////////////////////////////////////
/// Other
///////////////////////////////////////////////////////////////////////////////

func FindSolutionDirFrom(dir string) (string, error) {
  return findParentDirWithFile(sir, ".gosolution")
}

func HasSolution(dir string) bool {
  return isFile(fmt.Sprintf("%s/.gosolution", dir))
}
