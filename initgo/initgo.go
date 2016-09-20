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
	"reflect"
	"io/ioutil"
	"syscall"
	"bytes"
	"flag"
	"runtime/pprof"//性能分析
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
var programs []*program
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
	cmd := exec.Command("/bin/sh", "-c", `ps -ax | awk '{print $1}' | grep -e "^` + pid + `$"`)
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
func programSelf() {
	_, err := os.Stat(PidPath + ProgramName)
	if err == nil || os.IsExist(err) {
		btPid, err := ioutil.ReadFile(PidPath + ProgramName)
		if err == nil {
			cmd := exec.Command("/bin/sh", "-c", `"kill ` + string(btPid) + `"`)
			cmd.Start()
			time.Sleep(time.Second * 3)
		}
	}

	pid := os.Getpid()
	err = ioutil.WriteFile(PidPath + ProgramName, []byte(strconv.Itoa(pid)), 0660)
	if err != nil {
		panic(err)
		log.Error("%s", err)
		fmt.Println(err)
		os.Exit(0)
	}
	//CTRL+C退出
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-c
		programKill()
	}()
}

func programKill() {
	if SysCall != nil {
		SysCall()
	}
	if isStartCPUProfile {
		pprof.StopCPUProfile()
	}
	log.Info("%s", "程序结束成功")
	os.Remove(PidPath + ProgramName)
	os.Exit(0)
}
func InitProgram(Pgs ... Program) {
	if (programs != nil) {
		log.Alert("%s", "请勿重启调用")
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	log = logs.NewLogger(10000)
	log.EnableFuncCallDepth(true)
	if !Debug {
		log.SetLogger("file", `{"filename":"` + LogPath + `"}`)
	} else {
		log.SetLogger("console", "")
	}

	programSelf()

	if (len(Pgs) <= 0) {
		return
	}
	for _, v := range Pgs {
		if reflect.TypeOf(v).Name() == "Program" {
			ps := new(program)
			ps.pg = v
			programs = append(programs, ps)
		}
	}
}

func AddProgram(Pg Program) {
	if (Pg == nil || programs == nil) {
		return
	}
	if reflect.TypeOf(Pg).Name() == "Program" {
		ps := new(program)
		ps.pg = Pg
		programs = append(programs, ps)
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write mem profile to file")
var isStartCPUProfile bool

func SetCpuProfile(path string) {
	cpuprofile = &path
}

func SetMemProfile(path string) {
	memprofile = &path
}

//运行
func Run(Pgs ...Program) {
	if *cpuprofile == "" || *memprofile == "" {
		flag.Parse()
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Error("%s", err)
		} else {
			//统计cpu消耗
			pprof.StartCPUProfile(f)
			isStartCPUProfile = true

		}
	}
	if *memprofile != "" {
		fm, err := os.Create(*memprofile)
		if err != nil {
			log.Error("%s", err)
		} else {
			//统计内存消耗
			pprof.WriteHeapProfile(fm)
			defer fm.Close()
		}
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
func (this *program)virtual() {
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
