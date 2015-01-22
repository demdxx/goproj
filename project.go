// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package goproj

import (
  "errors"
  "fmt"
  "net/url"
  "path/filepath"
  "strings"

  "github.com/demdxx/gocast"
)

type Project struct {
  Dependency
  Owner interface{} // Project high level owner
  Deps  []*Dependency
}

func ProjectFromUrl(args ...string) *Project {
  p := &Project{}
  if IsUrl(args[0]) {
    p.Url = args[0]
    if len(args) > 1 {
      p.Path = args[1]
    } else {
      var err error
      p.Path, err = PathFromUrl(p.Url)
      if nil != err {
        return nil
      }
      p.Path = p.Path
    }
  } else {
    p.Path = args[0]
  }
  p.Name = p.Path
  p.Config = Config{"project": p.Name}
  return p
}

func ProjectFromFile(solpath, projpath string, conf Config) (proj *Project, err error) {
  name := projpath
  projConf := Config{"project": name}
  if err = projConf.InitFromFile(ProjectFileFullPathFor(solpath, projpath)); nil != err {
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
  if _, ok := proj.Config["version"]; !ok {
    proj.Config["version"] = "0"
  }

  if url, ok := proj.Config["url"]; ok {
    proj.Dependency.Url = gocast.ToString(url)
  }

  err = proj.Init()
  return
}

func (proj *Project) Init() error {
  deps, ok := proj.Config["deps"]
  if !ok {
    return nil
  }

  // Reset deps
  proj.Deps = nil

  // Init deps
  switch deps.(type) {
  case []interface{}:
    for _, depconf := range deps.([]interface{}) {
      switch depconf.(type) {
      case string:
        proj.addDependencyByConfig(depconf.(string), nil)
      default:
        if url, ok := depconf.(Config)["url"]; ok {
          proj.addDependencyByConfig(url.(string), depconf.(Config))
        } else {
          return errors.New(fmt.Sprintf("Invalid project[%s] config deps", proj.Path))
        }
      }
    }
  case map[string]interface{}:
    for url, depconf := range deps.(map[string]interface{}) {
      proj.addDependencyByConfig(url, Config(depconf.(map[string]interface{})))
    }
  }
  return nil
}

func (proj *Project) ReloadConfig() (err error) {
  projConf := Config{}
  if err = projConf.InitFromFile(proj.ProjectFullPath()); nil != err {
    return
  }

  proj.Config = projConf
  err = proj.Init()
  return
}

// Init FS struct
//
// .goproj
func (proj *Project) InitFileStruct(solpath string) error {
  if len(proj.Path) < 1 {
    if len(proj.Url) > 0 {
      path, err := ProjectPathFromUrl(proj.Url)
      if nil != err {
        return err
      }
      proj.Path = path
    } else {
      return errors.New("Solution path not defined")
    }
  }
  if err := makeDir(ProjectFullPathFor(solpath, proj.Path)); nil != err {
    return err
  }

  if len(proj.Url) > 0 {
    // Load project from URL
    _, cmd, url := PrepareCVSUrl(proj.Url)
    var err error
    var command interface{}

    if command, err = prepareCommand(proj, strings.Replace(cmd, "{url}", url, -1), nil); nil != err {
      return err
    }

    if err := run(proj, command.(string)); nil != err {
      return err
    }
    return proj.Init()
  }

  // Create solution
  return proj.SaveConfig()
}

// Save project
//
// @return nil or error
func (proj *Project) SaveConfig() error {
  return proj.Config.Save(proj.ProjectFileFullPath())
}

// TODO Init enviroment before run any command
func (proj *Project) UpdateEnv() {
  proj.Dependency.UpdateEnv()
}

///////////////////////////////////////////////////////////////////////////////
/// Actions
///////////////////////////////////////////////////////////////////////////////

func (proj *Project) CmdExec(cmd string, args []string, flags map[string]interface{}, observe bool) error {
  // Before run for dependencies
  if (nil == args || len(args) < 1) && nil != proj.Deps && len(proj.Deps) > 0 {
    for _, d := range proj.Deps {
      execute(d, cmd, flags, observe)
    }
  }

  // Run commands for me
  return execute(proj, cmd, flags, observe)
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
  dep := &Dependency{Name: url, Url: url, Config: config, Owner: proj}
  proj.Deps = append(proj.Deps, dep)
}

///////////////////////////////////////////////////////////////////////////////
/// Commands
///////////////////////////////////////////////////////////////////////////////

func (proj *Project) Cmds() map[string]interface{} {
  return proj.Dependency.Cmds()
}

func (proj *Project) Cmd(name string, def interface{}) interface{} {
  return proj.Dependency.Cmd(name, def)
}

func (proj *Project) CmdGet() interface{} {
  return nil // proj.Dependency.CmdGet()
}

// @return {go} build {flags} {app} or custom
func (proj *Project) CmdBuild() interface{} {
  return proj.Cmd("build", "{go} build {flags} {app}")
}

func (proj *Project) CmdInstall() interface{} {
  return proj.Cmd("install", "{go} install {flags} {app}")
}

func (proj *Project) CmdRun() interface{} {
  return proj.Dependency.CmdRun()
}

func (proj *Project) CmdTest() interface{} {
  return proj.Dependency.CmdTest()
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func (proj *Project) ProjectFullPath() string {
  return ProjectFullPathFor(proj.SolutionPath(), proj.Path)
}

func (proj *Project) ProjectFileFullPath() string {
  return ProjectFileFullPathFor(proj.SolutionPath(), proj.Path)
}

func ProjectFullPathFor(solpath, path string) string {
  if len(solpath) > 0 && !filepath.IsAbs(path) {
    return fmt.Sprintf("%s/src/%s/", solpath, path)
  }
  return path + "/"
}

func ProjectFileFullPathFor(solpath, path string) string {
  return ProjectFullPathFor(solpath, path) + "/.goproj"
}

///////////////////////////////////////////////////////////////////////////////
/// Other
///////////////////////////////////////////////////////////////////////////////

func ProjectPathFromUrl(u string) (string, error) {
  _url, err := url.Parse(u)
  if nil != err {
    return "", err
  }
  return fmt.Sprintf("%s%s", _url.Host, _url.Path), nil
}

func FindProjectDirFrom(dir string) (string, error) {
  return findParentDirWithFile(dir, ".goproj")
}
