package goproj

type Project Dependency

// @return {go} build {flags} {app} or ""
func (d Dependency) BuildCmd() string {
  if cmd, ok := d.Config["build"]; ok {
    return cmd
  }
  return "{go} build {flags} {app}"
}
