// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "errors"
  "fmt"
  "path/filepath"
)

type Project struct {
  Dependency
  Deps []*Dependency
}

func ProjectFromFile(solpath, projpath string, conf Config) (proj *Project, err error) {
  name := projpath
  if !filepath.IsAbs(projpath) {
    projpath, err = filepath.Abs(solpath + "/src/" + projpath)
    if nil != err {
      return
    }
  }

  projConf := Config{}
  if err = projConf.InitFromFile(projpath + "/.goproj"); nil != err {
    return
  }

  // Preinit project
  proj = &Project{
    Dependency: Dependency{
      Name:   name,
      Path:   projpath,
      Config: projConf,
    },
  }

  // Merge local and global config
  proj.Config.Update(conf, false)
  err = proj.Init()
  return
}

func (proj *Project) Init() error {
  deps, ok := proj.Config["deps"]
  if !ok {
    return errors.New("Project don't have dependencies")
  }

  // Reset deps
  proj.Deps = nil

  // Init deps
  switch deps.(type) {
  case []interface{}:
    for _, depconf := range deps.([]interface{}) {
      if url, ok := depconf.(Config)["url"]; ok {
        proj.addDependencyByConfig(url.(string), depconf.(Config))
      } else {
        return errors.New(fmt.Sprintf("Invalid project[%s] config deps", proj.Path))
      }
    }
  case map[string]interface{}:
    for url, depconf := range deps.(map[string]interface{}) {
      proj.addDependencyByConfig(url, Config(depconf.(map[string]interface{})))
    }
  }
  return nil
}

// Init FS struct
//
// .goproj
func (proj *Project) InitFileStruct() error {
  if len(proj.Path) < 1 {
    return errors.New("Solution path not defined")
  }
  if err := makeDir(proj.Path); nil != err {
    return err
  }

  // Create solution
  return proj.SaveConfig()
}

// Save project
//
// @return nil or error
func (proj *Project) SaveConfig() error {
  return proj.Config.Save(fmt.Sprintf("%s/.goproj", proj.Path))
}

// TODO Init enviroment before run any command
func (proj *Project) UpdateEnv() {
  proj.Dependency.UpdateEnv()
}

///////////////////////////////////////////////////////////////////////////////
/// Actions
///////////////////////////////////////////////////////////////////////////////

func (proj *Project) CmdExec(cmd string, args []string, flags map[string]interface{}) error {
  // Before run for dependencies
  if (nil == args || len(args) < 1) && nil != proj.Deps && len(proj.Deps) > 0 {
    for _, d := range proj.Deps {
      execute(d, cmd, flags)
    }
  }

  // Run commands for me
  return execute(proj, cmd, flags)
}

///////////////////////////////////////////////////////////////////////////////
/// Dependencies
///////////////////////////////////////////////////////////////////////////////

// Get dependency by index
func (proj Project) DependencyByIndex(index int) (dep *Dependency, err error) {
  if nil == proj.Deps || len(proj.Deps) <= index {
    err = errors.New(fmt.Sprintf("Undefined dependency by index %d", index))
  } else {
    dep = proj.Deps[index]
  }
  return
}

// Get dependency by url
func (proj Project) DependencyByUrl(url string) (dep *Dependency, err error) {
  if nil == proj.Deps {
    err = errors.New("Project don't have dependencies")
    return
  }
  for _, depend := range proj.Deps {
    if url == depend.Url {
      dep = depend
      return
    }
  }
  err = errors.New(fmt.Sprintf("Undefined dependency by url: %s", url))
  return
}

// Add dependency to deps list
// @param url
// @param conf
func (proj *Project) addDependencyByConfig(url string, conf Config) {
  if nil == proj.Deps {
    proj.Deps = make([]*Dependency, 0)
  }

  // Update config
  config := Config{}
  config.Update(conf, true)

  // Init project
  dep := &Dependency{Url: url, Config: config}
  proj.Deps = append(proj.Deps, dep)
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func (proj *Project) Cmds() map[string]interface{} {
  return proj.Dependency.Cmds()
}

func (proj *Project) Cmd(name string, def interface{}) interface{} {
  return proj.Dependency.Cmd(name, def)
}

func (proj *Project) CmdGet() interface{} {
  return proj.Dependency.CmdGet()
}

// @return {go} build {flags} {app} or custom
func (proj *Project) CmdBuild() interface{} {
  return proj.Cmd("build", "{go} build {flags} {app}")
}

func (proj *Project) CmdRun() interface{} {
  return proj.Dependency.CmdRun()
}