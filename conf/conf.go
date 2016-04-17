package conf

import (
  "os"
  pathPackage "path"
  "strings"
  "errors"
)

var (
  ConfDirName = "conf"
  ConfNamePrefix = "app"
  ConfSuffix = ".ini"
  ConfCycles = 5
  confTotal = 0;
)
func FindConf(path string) (string,error) {

  if confTotal > ConfCycles {
    return "",errors.New("cycles Upper limit")
  }
  info, err := os.Lstat(path)
  if err != nil {
    return "",err
  }

  if !info.IsDir() {
    return "",errors.New("not dir")
  }


  names, err := readDirNames(path)
  if err != nil {
    return "",err
  }
  var confTmp string
  for _, name := range names {
    //fmt.Println(name)

    //特征相同
    if strings.HasPrefix(name, ConfNamePrefix) && strings.HasSuffix(name, ConfSuffix) {
      confTmp = pathPackage.Join(path, name)
      f, err := os.Stat(confTmp)
      if err == nil && !f.IsDir() {
        return confTmp,nil

      }
    }else if name == ConfDirName {
      confTmp = pathPackage.Join(path, name)
      f, err := os.Stat(confTmp)
      if err == nil && f.IsDir() {
        return FindConf(confTmp)
      }
    }
  }
  confTmp = pathPackage.Join(path,"../")

  confTotal++

  return FindConf(confTmp)

}
func readDirNames(dirname string) ([]string, error) {
  f, err := os.Open(dirname)
  if err != nil {
    return nil, err
  }
  names, err := f.Readdirnames(-1)
  f.Close()
  if err != nil {
    return nil, err
  }
  return names, nil
}