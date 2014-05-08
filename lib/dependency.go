package lib

type Dependency struct {
  Path   string
  Url    string
  Config Config
}

// @return {cmd} or ""
func (d Dependency) BuildCmd() string {
  if cmd, ok := d.Config["build"]; ok {
    return cmd.(string)
  }
  return ""
}
