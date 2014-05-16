package lib

func indexOfStringSlice(slice []string, s string) int {
  if nil != slice {
    for i, v := range slice {
      if s == v {
        return i
      }
    }
  }
  return -1
}
