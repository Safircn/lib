package rand

import (
  "time"
  "math/rand"
)


var randString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func GetRandString(length int) string{
  bytes := []byte(randString)
  result := []byte{}
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for i := 0; i < length; i++ {
    result = append(result, bytes[r.Intn(len(bytes))])
  }
  return string(result)
}
