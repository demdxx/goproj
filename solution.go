package goproj

import (
  "bytes"
  "encoding/json"
  "gopkg.in/v1/yaml"
  "os"
)

type Solution struct {
  Path     string
  Projects []*Project
  Config   Config
}

func SolutionFromFile(filepath string) (sol *Solution, err error) {
  file, err := os.Open(filepath)
  if err != nil {
    return
  }
  defer file.Close()

  // Load file
  reader := bufio.NewReader(file)
  buffer := bytes.NewBuffer()
  for {
    if part, prefix, err = reader.ReadLine(); err != nil {
      break
    }
    buffer.Write(part)
  }
  if err == io.EOF {
    err = nil
  }

  // Load solution config
  sol = &Solution{}
  data := buffer.Bytes()
  if '{' == data[0] {
    err = json.Unmarshal(data, &sol.Config)
  } else {
    err = yuml.Unmarshal(data, &sol.Config)
  }

  // Init project
  if nil == err {
    err = sol.Init()
  }
  return
}

func (sol Solution) Init() error {
  if nil == sol.Config || len(sol.Config) < 1 {
    return errors.New("Project not inited")
  }

  if nil != sol.Projects {
    sol.Projects = nil
  }

  // Each config
  for dir, conf := range sol.Config {
    proj := ProjectFromFile()
  }

  return nil
}

func (sol Solution) AddProject(p *Project) error {
  if nil == sol.Projects {
    sol.Projects = make([]*Project, 0)
  }
  sol.Projects = append(sol.Projects, p)
  return nil
}
