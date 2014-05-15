package lib

import (
  "net/url"
  "strings"
)

var CVS_tpl = map[string]string{
  "git": "git clone {url} {path}",
  "hg":  "hg clone {url} {path}",
  "bzr": "bzr branch {url} {path}",
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
      cvs = u.User.Username()
      if cmd, ok = CVS_tpl[cvs]; ok {
        return
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
