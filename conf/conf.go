package conf

import (
  "os"
  pathPackage "path"
  "strings"
  "errors"
  "sort"
)

var (
  ConfDirName                                     = "conf"
  ConfNamePrefix                                  = "app."
  ConfSuffix                                      = ".ini"
  ConfCycles                                      = 5
  confTotal                                       = 0;
  ConfPathsSelectFunc func(paths []string) string = func(paths []string) string {
    if len(paths) == 1 {
      return paths[0]
    }
    sort.Strings(paths)
    return paths[0]
  }
)

func FindConf(path string) (string, error) {
  paths := make([]string, 0)
  err := findConf(path, &paths, false)
  if err != nil {
    return "", err
  }
  if len(paths) == 0 {
    return "", errors.New("not dir")
  }
  path = ConfPathsSelectFunc(paths)
  if path == "" {
    return "", errors.New("not dir")
  }
  return path, nil
}

func findConf(path string, paths *[]string, isConfDirNameDir bool) (error) {
  if confTotal > ConfCycles {
    return errors.New("cycles Upper limit")
  }
  info, err := os.Lstat(path)
  if err != nil {
    return err
  }

  if !info.IsDir() {
    return errors.New("not dir")
  }
  names, err := readDirNames(path)
  if err != nil {
    return err
  }
  var confTmp string
  for _, name := range names {
    //特征相同
    if isConfDirNameDir {
      if strings.HasPrefix(name, ConfNamePrefix) && strings.HasSuffix(name, ConfSuffix) {
        confTmp = pathPackage.Join(path, name)
        f, err := os.Stat(confTmp)
        if err == nil && !f.IsDir() {
          *paths = append(*paths, confTmp)
        }
      }
    } else if name == ConfDirName {
      confTmp = pathPackage.Join(path, name)
      f, err := os.Stat(confTmp)
      if err == nil && f.IsDir() {
        return findConf(confTmp, paths, true)
      }
    }
  }
  if len(*paths) > 0 {
    return nil
  }
  confTmp = pathPackage.Join(path, "../")

  confTotal++

  return findConf(confTmp, paths, false)
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
