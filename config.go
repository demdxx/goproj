package goproj

type Config map[string]interface{}

func (conf Config) InitFromFile(filepath string) (err error) {
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
  if nil == err {
    // Load solution config
    data := buffer.Bytes()
    if '{' == data[0] {
      err = json.Unmarshal(data, &conf)
    } else {
      err = yuml.Unmarshal(data, &conf)
    }
  }
  return
}
