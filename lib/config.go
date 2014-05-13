// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "bufio"
  "bytes"
  "encoding/json"
  "gopkg.in/v1/yaml"
  "io"
  "os"
)

type Config map[string]interface{}

func (conf *Config) InitFromFile(filepath string) (err error) {
  file, err := os.Open(filepath)
  if err != nil {
    return
  }
  defer file.Close()

  // Load file
  reader := bufio.NewReader(file)
  var buffer bytes.Buffer
  var part []byte
  for {
    if part, _, err = reader.ReadLine(); err != nil {
      break
    }
    buffer.Write(part)
  }
  if err == io.EOF {
    err = nil
  }
  if nil == err {
    // Load solution config
    data := buffer.Bytes()
    if '{' == data[0] {
      err = json.Unmarshal(data, &conf)
    } else {
      err = yaml.Unmarshal(data, &conf)
    }
  }
  return
}

func (conf *Config) Update(c Config, rewrite bool) {
  if nil != conf {
    for k, v := range conf {
      if rewrite {
        proj.Config[k] = v
      } else if _, ok := proj.Config[k]; !ok {
        proj.Config[k] = v
      }
    }
  }
}
