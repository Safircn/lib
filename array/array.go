package array

func InArray(s string,arr []string) bool {
  for _,v := range arr {
    if v == s {
      return true
    }
  }
  return false
}
