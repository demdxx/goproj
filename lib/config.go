// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package lib

import (
  "encoding/json"
  "gopkg.in/v1/yaml"
  "io/ioutil"
)

type Config map[string]interface{}

func (conf Config) InitFromFile(filepath string) (err error) {
  // Load file
  var data []byte
  data, err = ioutil.ReadFile(filepath)

  if nil == err {
    // Parse solution config
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
func (conf Config) Save(fpath string) error {
  data, err := yaml.Marshal(conf)
  if nil != err {
    return err
  }
  return ioutil.WriteFile(fpath, data, 0644)
}

func ConfigConvert(data interface{}) Config {
  switch data.(type) {
  case map[string]interface{}:
    return data.(Config)
    break
  case map[interface{}]interface{}:
    conf := make(Config)
    for k, v := range data.(map[interface{}]interface{}) {
      conf[k.(string)] = v
    }
    return conf
    break
  }
  return nil
}
