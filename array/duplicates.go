package array

func RemoveDuplicatesAndEmptyStr(list []string) []string {
  var (
    x []string = []string{}
    flag bool
  )
  for _, i := range list {
    if len(x) == 0 {
      x = append(x, i)
    } else {
      flag = false
      for _, v := range x {
        if i == v {
          flag = true
          break
        }
      }
      if !flag {
        x = append(x, i)
      }
    }
  }
  return x
}


func RemoveDuplicatesAndEmptyInt(list []int) []int {
  var (
    x []int = []int{}
    flag bool
  )
  for _, i := range list {
    if len(x) == 0 {
      x = append(x, i)
    } else {
      flag = false
      for _, v := range x {
        if i == v {
          flag = true
          break
        }
      }
      if !flag {
        x = append(x, i)
      }
    }
  }
  return x
}

func RemoveDuplicatesAndEmptyInt64(list []int64) []int64 {
  var (
    x []int64 = []int64{}
    flag bool
  )
  for _, i := range list {
    if len(x) == 0 {
      x = append(x, i)
    } else {
      flag = false
      for _, v := range x {
        if i == v {
          flag = true
          break
        }
      }
      if !flag {
        x = append(x, i)
      }
    }
  }
  return x
}
