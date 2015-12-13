// @project goproj
// @copyright Dmitry Ponomarev <demdxx@gmail.com> 2014
//
// This work is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/.

package goproj

import (
  "encoding/json"
  "gopkg.in/v1/yaml"
  "io/ioutil"

  "github.com/demdxx/gocast"
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

func (conf Config) GetString(key, def string) string {
  if v, ok := conf[key]; ok && nil != v {
    return gocast.ToString(v)
  }
  return def
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func ConfigConvert(data interface{}) Config {
  switch data.(type) {
  case map[string]interface{}:
    conf := make(Config)
    for k, v := range data.(map[string]interface{}) {
      conf[k] = v
    }
    return conf
  case map[interface{}]interface{}:
    conf := make(Config)
    data, _ = gocast.ToSiMap(data, "", false)
    if nil != data {
      for k, v := range data.(map[string]interface{}) {
        conf[k] = v
      }
    }
    return conf
  }
  return nil
}
