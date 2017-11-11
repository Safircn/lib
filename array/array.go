package array

func InArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func InArrayStr(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func InArrayInt(s int, arr []int) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func InArrayInt64(s int64, arr []int64) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func Rm_duplicate(list []int) []int {
  var x []int = []int{}
  var flag bool
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