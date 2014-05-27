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
  "strings"
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
    pdir, err := findParentDirWithFile(dir, ".goproj")
    if nil != err {
      return nil, err
    }

    // Init global solution
    GOPATH := os.Getenv("GOPATH")
    if len(GOPATH) < 1 {
      err = errors.New("For global solution project need set enviroment GOPATH")
    } else {
      err = nil
      sol = new(Solution)
      sol.IsGlobal = true
      sol.Config["projects"] = map[string]interface{}{
        pdir: map[string]interface{}{}, // Empty config
      }
      return sol, sol.Init(GOPATH)
    }
  } else {
    // Init solution from file
    sol, err = SolutionFromFile(fmt.Sprintf("%s/.gosolution", ndir))
  }
  return
}

func SolutionFromFile(fpath string) (sol *Solution, err error) {
  conf := Config{}
  if err = conf.InitFromFile(fpath); nil == err {
    sol = new(Solution)
    sol.Config = conf
    err = sol.Init(filepath.Dir(fpath))
  }
  return
}

func (sol *Solution) Init(path string) (err error) {
  if len(path) > 0 {
    sol.Path, err = filepath.Abs(path)
    if nil != err {
      return
    }
  }

  if nil == sol.Config || len(sol.Config) < 1 {
    err = errors.New("Project not inited")
    return
  }

  if nil != sol.Projects {
    sol.Projects = nil
  }

  // Each config
  if projects, ok := sol.Config["projects"]; ok {
    if nil != projects {
      switch projects.(type) {
      case map[interface{}]interface{}:
      case map[string]interface{}:
        projs := ToStringMap(projects)
        for dir, conf := range projs {
          var proj *Project
          proj, err = ProjectFromFile(sol.Path, dir, ConfigConvert(conf))
          if nil == err && nil != proj {
            err = sol.AddProject(proj)
          }
          if nil != err {
            return
          }
        }
        break
      default:
        err = errors.New("Config has invalid format in conf.projects section")
        return
      }
    }
  }
  return
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
  if err := makeDir(fmt.Sprintf("%s/pkg", sol.Path)); nil != err {
    return err
  }
  if err := makeDir(fmt.Sprintf("%s/src", sol.Path)); nil != err {
    return err
  }

  // Process initialize projects
  if nil != sol.Projects && len(sol.Projects) > 0 {
    for _, p := range sol.Projects {
      if err := p.InitFileStruct(sol.Path + "/src/"); nil != err {
        return err
      }
    }
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
  return ioutil.WriteFile(fmt.Sprintf("%s/.gosolution", sol.Path), data, 0644)
}

///////////////////////////////////////////////////////////////////////////////
/// Actions
///////////////////////////////////////////////////////////////////////////////

func (sol *Solution) CmdExec(cmd string, args []string, flags map[string]interface{}) error {
  if nil != sol.Projects && len(sol.Projects) > 0 {
    // Process command
    for _, p := range sol.Projects {
      if nil == args || len(args) < 1 || -1 != indexOfStringSlice(args, p.Name) {
        // Init environment
        sol.UpdateEnv()

        // Do exec
        if err := p.CmdExec(cmd, args, flags); nil != err {
          return err
        }
      }
    }
  }
  return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Getters/Setters
///////////////////////////////////////////////////////////////////////////////

func (sol *Solution) EnvPath() string {
  PATH := os.Getenv("PATH")
  return fmt.Sprintf("%s/bin:%s", strings.TrimRight(sol.Path, "/"), PATH)
}

func (sol *Solution) UpdateEnv() {
  if !sol.IsGlobal {
    os.Setenv("GOPATH", sol.Path)
    os.Setenv("PATH", sol.EnvPath())
  }
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
  p.Owner = sol
  sol.Projects = append(sol.Projects, p)
  return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Other
///////////////////////////////////////////////////////////////////////////////

func FindSolutionDirFrom(dir string) (string, error) {
  return findParentDirWithFile(dir, ".gosolution")
}

func HasSolution(dir string) bool {
  b, _ := isFile(fmt.Sprintf("%s/.gosolution", dir))
  return b
}

// Get enviroment for solution or project dir
//
// @param dir path
// @return map{GOPATH,PATH,GO}
func SolutionEnv(dir string) map[string]string {
  var GOPATH, PATH string
  PATH = os.Getenv("PATH")

  sol, _ := SolutionFromDir(dir)
  if nil != sol && !sol.IsGlobal {
    GOPATH = sol.Path
    PATH = sol.EnvPath()
  } else {
    GOPATH = os.Getenv("GOPATH")
    PATH = fmt.Sprintf("%s/bin:%s", strings.TrimRight(GOPATH, "/"), PATH)
  }
  return map[string]string{
    "GOPATH": GOPATH,
    "PATH":   PATH,
    "GO":     GoPath(),
  }
}

///////////////////////////////////////////////////////////////////////////////
/// Extends
///////////////////////////////////////////////////////////////////////////////

func (proj *Project) SolutionPath() string {
  if nil != proj.Owner {
    return proj.Owner.(*Solution).Path
  }
  return ""
}

func (proj Project) Solution() *Solution {
  return proj.Owner.(*Solution)
}

func (dep *Dependency) SolutionPath() string {
  if nil != dep.Owner {
    return dep.Owner.(*Project).SolutionPath()
  }
  return ""
}

func (dep Dependency) Solution() *Solution {
  return dep.Owner.(*Project).Solution()
}
