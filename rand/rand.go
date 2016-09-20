package rand

import (
  "time"
  "math/rand"
)

const (
  RAND_KIND_NUM   = 0  // 纯数字
  RAND_KIND_LOWER = 1  // 小写字母
  RAND_KIND_UPPER = 2  // 大写字母
  RAND_KIND_ALL   = 3  // 数字、大小写字母
)

func GetRandString(length int) string{
  return string(Krand(length,RAND_KIND_ALL))
}


// 随机字符串
func Krand(size int, kind int) []byte {
  ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
  rand.Seed(time.Now().UnixNano())
  for i :=0; i < size; i++ {
    if kind == 3 { // random ikind
      ikind = rand.Intn(3)
    }
    scope, base := kinds[ikind][0], kinds[ikind][1]
    result[i] = byte(base+rand.Intn(scope))
  }
  return result
}