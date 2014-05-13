// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

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
