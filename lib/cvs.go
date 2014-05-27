// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "net/url"
  "strings"
)

var CVS_tpl = map[string]string{
  "git": "git clone {url} {fullpath}",
  "hg":  "hg clone {url} {fullpath}",
  "bzr": "bzr branch {url} {fullpath}",
  "go":  "go get {flags} {url}",
}

func PrepareCVSUrl(_url string) (cvs, cmd, p_url string) {
  var ok bool

  // set prepared url
  p_url = _url

  if len(_url) > 7 {
    parts := strings.Split(_url, ":")
    if len(parts) > 1 {
      cvs = parts[0]
      if cmd, ok = CVS_tpl[cvs]; ok {
        nurl := _url[len(cvs)+1:]
        if "//" != nurl[:2] {
          p_url = nurl
          return
        }
      }
    }

    // Check url {cvs}@{url}.{cvs}
    u, err := url.Parse(_url)
    if nil == err {
      // Check user name
      if nil != u.User {
        cvs = u.User.Username()
        if cmd, ok = CVS_tpl[cvs]; ok {
          return
        }
      }

      // Check path tail
      parts = strings.Split(u.Path, ".")
      cvs = parts[len(parts)-1]
      if cmd, ok = CVS_tpl[cvs]; ok {
        return
      }
    }
  }
  cvs, cmd = "go", CVS_tpl["go"]
  return
}
