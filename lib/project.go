package lib

import (
  "errors"
  "fmt"
)

type Project struct {
  Dependency
  Deps []*Dependency
}

func ProjectFromFile(filepath string) (proj *Project, err error) {
  proj = &Project{}
  if err = proj.Config.InitFromFile(filepath); nil != err {
    proj = nil
    return
  }
  err = proj.Init()
  return
}

func (proj Project) Init() error {
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
      proj.addDependencyByConfig(url, depconf.(Config))
    }
  }
  return nil
}

// @return {go} build {flags} {app} or ""
func (proj Project) BuildCmd() string {
  if cmd, ok := proj.Config["build"]; ok {
    return cmd.(string)
  }
  return "{go} build {flags} {app}"
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
func (proj Project) addDependencyByConfig(url string, conf Config) {
  if nil == proj.Deps {
    proj.Deps = make([]*Dependency, 0)
  }
  dep := &Dependency{Url: url, Config: conf}
  proj.Deps = append(proj.Deps, dep)
}
