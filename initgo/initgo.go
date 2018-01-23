package initgo

import (
  "runtime"
  "fmt"
  "github.com/Safircn/lib/logs"
  "strconv"
  "time"
  "os"
  "os/signal"
  "os/exec"
  "io/ioutil"
  "syscall"
  "bytes"
  "errors"
)

type Program func()

var ProgramStatusNum int64
var log *logs.BeeLogger

type program struct {
  pg         Program
  status     int8      //当前程序状态 0为尚未开启 1开启 2程序意外挂死 3运行结束
  startTime  time.Time //启动时间戳
  restartNum int64     //重启次数
  endTime    time.Time //结束时间
}

//方法列表列表
var programs []*program = make([]*program, 0)
var ProgramName = "golang.pid";

var PidPath = "/var/run/"
var Debug = false
var LogPath = "initVirtual.log"
var SysCall func()

/**
进程是否运行
 */
func IsRun() bool {
  pid, err := GetPid()
  if err != nil {
    return false
  }
  //✗ ps -ax |awk '{print $1}' | grep -e "^94811$"
  cmd := exec.Command("/bin/sh", "-c", `ps -ax | awk '{print $1}' | grep -e "`+pid+`"`)
  var outBytes []byte
  outBytes, err = cmd.Output()
  if err == nil && bytes.Equal(bytes.TrimSpace(outBytes), []byte(pid)) {
    return true
  }
  return false
}

func GetPid() (pid string, err error) {
  _, err = os.Stat(PidPath + ProgramName)
  if err == nil || os.IsExist(err) {
    var btPid []byte
    btPid, err = ioutil.ReadFile(PidPath + ProgramName)
    if err == nil {
      pid = string(bytes.TrimSpace(btPid))
      return
    }
    return
  } else {
    return
  }
}

//程序成活处理  pid 处理
func programSelf() error {
  _, err := os.Stat(PidPath + ProgramName)
  if err == nil || os.IsExist(err) {
    btPid, err := ioutil.ReadFile(PidPath + ProgramName)
    if err == nil {
      cmds := exec.Command("/bin/sh", "-c", fmt.Sprintf(`ps -ax | awk '{print $1}' | grep -e "%s"`, btPid))
      outptBytes, err := cmds.Output()

      if err != nil {
        if serr, ok := err.(*exec.ExitError); ok {
          if waitStatus, ok2 := serr.ProcessState.Sys().(syscall.WaitStatus); ok2 {
            if waitStatus.ExitStatus() == 1 {
              goto PROGRAM_SELF_NEXT
            }
          }
        }
        return errors.New("ps error command:" + fmt.Sprintf(`ps -ax | awk '{print $1}' | grep -e "%s"`, btPid) + " err:" + err.Error())
      }
    PROGRAM_SELF_NEXT:
      if len(outptBytes) > 0 {
        return errors.New("pid already exist pid:" + string(btPid))
      } else {
        err = os.Remove(PidPath + ProgramName)
        if err != nil {
          return errors.New("pidFile Remove err:" + err.Error())
        }
      }
    }
  }

  pid := os.Getpid()
  err = ioutil.WriteFile(PidPath+ProgramName, []byte(strconv.Itoa(pid)), 0660)
  if err != nil {
    return err
  }
  //CTRL+C退出
  go func() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
    <-c
    programKill()
  }()
  return nil
}

func programKill() {
  if SysCall != nil {
    SysCall()
  }
  log.Info("%s", "程序结束成功")
  os.Remove(PidPath + ProgramName)
  os.Exit(0)
}

func InitProgram(Pgs ... Program) {
  if (len(Pgs) == 0) {
    return
  }
  for _, v := range Pgs {
    ps := new(program)
    ps.pg = v
    programs = append(programs, ps)
  }
}

func AddProgram(Pg Program) {
  if (Pg == nil) {
    return
  }
  ps := new(program)
  ps.pg = Pg
  programs = append(programs, ps)
}

//运行
func Run(Pgs ...Program) {
  runtime.GOMAXPROCS(runtime.NumCPU())
  serInit()

  err := programSelf()
  if err != nil {
    log.Error("init err:%s", err.Error())
    return
  }

  //执行初始化
  if (len(Pgs) > 0 ) {
    for _, pg := range Pgs {
      pg()
    }
  }
  if programs == nil {
    return
  }
  log.Info("%s", "守护程序启动中")
  for {
    ProgramStatusNum = 0
    for _, v := range programs {
      switch v.status {
      case 0:
        v.status = 1
        v.startTime = time.Now()
        go v.virtual()
        log.Info("%d:%s", ProgramStatusNum, "启动完成")
        ProgramStatusNum++
      case 1:
        ProgramStatusNum++
      case 2:
        v.restartNum++
        log.Info("程序重新启动:" + strconv.FormatInt(v.restartNum, 10))
        v.status = 1
        go v.virtual()
        ProgramStatusNum++
      }
    }
    if ProgramStatusNum == 0 {
      programKill()

    }
    time.Sleep(time.Second * 10)
  }
}
func serInit() {

  log = logs.NewLogger(10000)
  log.EnableFuncCallDepth(true)
  if !Debug {
    log.SetLogger("file", `{"filename":"`+LogPath+`"}`)
  } else {
    log.SetLogger("console", "")
  }
}

//捕获错误
/**
defer initgo.RecoverLog
 */
func RecoverLog() {
  if err := recover(); err != nil {
    log.EnableFuncCallDepth(false)
    log.Error("%s", err)
    for i := 1; ; i++ {
      _, file, line, ok := runtime.Caller(i)
      if !ok {
        break
      }
      log.Error(fmt.Sprintf("%s:%d", file, line))
    }
    log.EnableFuncCallDepth(true)
  }
}

//虚拟
func (this *program) virtual() {
  defer func() {
    if err := recover(); err != nil {

      log.EnableFuncCallDepth(false)
      log.Error("%s", err)
      for i := 1; ; i++ {
        _, file, line, ok := runtime.Caller(i)
        if !ok {
          break
        }
        log.Error(fmt.Sprintf("%s:%d", file, line))
        this.status = 2
      }
      log.EnableFuncCallDepth(true)
    }
  }()
  this.pg()
  this.status = 3
  this.endTime = time.Now()
  return
}
