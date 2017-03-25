package logger

/*
定制修改来自 github.com/donnie4w/go-logger/logger/logger.go
*/
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	_VER string = "1.0.3"
)

type LEVEL int32

var logLevel LEVEL = 1
var maxFileSize int64
var maxFileCount int32
var dailyRolling bool = true
var consoleAppender bool = true
var RollingFile bool = false
var logObj *_FILE

const DATEFORMAT = "2006-01-02"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	INFO
	WARN
	ERROR
	OFF
)

type _FILE struct {
	dir           string
	_suffix       int
	isCover       bool
	_date         *time.Time
	mu            *sync.RWMutex
	logfile_info  *os.File
	lg_info       *log.Logger
	logfile_warn  *os.File
	lg_warn       *log.Logger
	logfile_error *os.File
	lg_error      *log.Logger
}

func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

func SetLevel(_level LEVEL) {
	logLevel = _level
}

func LogHandler(fileDir string, num int) {

	for {
		DeleteLog(fileDir, num)
		time.Sleep(10 * time.Second)
	}

}

func DeleteLog(fileDir string, num int) {
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	dir_list, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("read dir:", fileDir, "error: ", err)
		return
	}

	for i, v := range dir_list {
		name := fileDir + "/" + v.Name()
		if v.IsDir() {
			DeleteLog(name, num)
			continue
		}
		if i < len(dir_list)-num {
			if isExist(name) {
				err := os.Remove(name)
				fmt.Println("os.Remove error:", err)
			}

		}

	}
}

func SetRollingDaily(fileDir string, num int) {
	RollingFile = false
	dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	logObj = &_FILE{dir: fileDir, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	go LogHandler(fileDir, num)
	createNewFile()
}
func createNewFile() {
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if logObj.logfile_info != nil {
		logObj.logfile_info.Close()
	}
	err := os.Chdir(logObj.dir + "/" + "info/")
	if err != nil {
		fmt.Println("os.Chdir(logObj.dir + / + info/)", err)
		err = os.Mkdir(logObj.dir+"/"+"info/", 0666)
		if err != nil {
			fmt.Println("os.Mkdir eror:", err)
		}
	}
	logObj.logfile_info, _ = os.OpenFile(logObj.dir+"/"+"info/"+logObj._date.Format(DATEFORMAT)+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	logObj.lg_info = log.New(logObj.logfile_info, "\n", log.Ldate|log.Ltime|log.Lshortfile)

	if logObj.logfile_warn != nil {
		logObj.logfile_warn.Close()
	}
	err = os.Chdir(logObj.dir + "/" + "warn/")
	if err != nil {
		fmt.Println("os.Chdir(logObj.dir + / + warn/)", err)
		err = os.Mkdir(logObj.dir+"/"+"warn/", 0666)
		if err != nil {
			fmt.Println("os.Mkdir eror:", err)
		}
	}
	logObj.logfile_warn, _ = os.OpenFile(logObj.dir+"/"+"warn/"+logObj._date.Format(DATEFORMAT)+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	logObj.lg_warn = log.New(logObj.logfile_warn, "\n", log.Ldate|log.Ltime|log.Lshortfile)

	if logObj.logfile_error != nil {
		logObj.logfile_error.Close()
	}
	err = os.Chdir(logObj.dir + "/" + "error/")
	if err != nil {
		fmt.Println("os.Chdir(logObj.dir + / + error/)", err)
		err = os.Mkdir(logObj.dir+"/"+"error/", 0666)
		if err != nil {
			fmt.Println("os.Mkdir eror:", err)
		}
	}
	logObj.logfile_error, _ = os.OpenFile(logObj.dir+"/"+"error/"+logObj._date.Format(DATEFORMAT)+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	logObj.lg_error = log.New(logObj.logfile_error, "\n", log.Ldate|log.Ltime|log.Lshortfile)
}
func console(s ...interface{}) {
	if consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(file+","+strconv.Itoa(line), s)
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println(fmt.Sprintf("catchError err:%v，logObj：%v", err, logObj))
	}
}

// func LogNetWarn(r *http.Request, v ...interface{}) {
// 	LogWarn(GetIPAddress(r), v)
// }
func LogInfo(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if logLevel <= INFO {
		if logObj != nil {
			err := logObj.lg_info.Output(2, fmt.Sprintln("info", v))
			if err != nil {
				fmt.Println("LogInfo output eror:", err)
			}
		}
		console("info", v)
	}
}
func LogWarn(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if logLevel <= WARN {
		if logObj != nil {
			err := logObj.lg_warn.Output(2, fmt.Sprintln("warn", v))
			if err != nil {
				fmt.Println("LogWarn output warn:", err)
			}
		}
		console("warn", v)
	}
}
func LogError(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if logLevel <= ERROR {
		if logObj != nil {
			err := logObj.lg_error.Output(2, fmt.Sprintln("error", v))
			if err != nil {
				fmt.Println("LogError output eror:", err)
			}
		}
		console("error", v)
	}
}

func (f *_FILE) isMustRename() bool {
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	if t.After(*f._date) {
		return true
	}
	return false
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(maxFileCount) + 1)
}

func fileSize(file string) int64 {
	fmt.Println("fileSize", file)
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			fileCheck()
		}
	}
}

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("fileCheck err:%v", err))
		}
	}()
	if logObj != nil && logObj.isMustRename() {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		logObj._date = &t
		createNewFile()
	}
}
