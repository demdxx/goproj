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
  "io/ioutil"
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

func (conf Config) Update(c Config, rewrite bool) {
  if nil != c {
    for k, v := range c {
      if rewrite {
        conf[k] = v
      } else if _, ok := conf[k]; !ok {
        conf[k] = v
      }
    }
  }
}

// Save project
//
// @return nil or error
func (conf *Config) Save(fpath string) error {
  data, err := yaml.Marshal(*conf)
  if nil != err {
    return err
  }
  return ioutil.WriteFile(fpath, data, 0644)
}
