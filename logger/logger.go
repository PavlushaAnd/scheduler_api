package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"scheduler_api/utils"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FileHdl       *os.File
)

func init() {
	DoInit()
}

func DoInit() {
	file, err := tryOpenLogfile()
	if err != nil {
		log.Fatal(err)
	}
	FileHdl = file
	DebugLogger = log.New(FileHdl, "D:", log.Ldate|log.Ltime)
	InfoLogger = log.New(FileHdl, "I:", log.Ldate|log.Ltime)
	WarningLogger = log.New(FileHdl, "W:", log.Ldate|log.Ltime)
	ErrorLogger = log.New(FileHdl, "E:", log.Ldate|log.Ltime)
}

func tryOpenLogfile() (file *os.File, err error) {
	logfilePath := getLogFileName()
	file, err = os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	return file, err
}

func getLogFileName() string {
	t := time.Now()
	filename := t.Format("20060102")
	filename += ".log"
	path, err := beego.AppConfig.String("LoggerRoot")
	if err != nil {
		log.Fatal(err)
	}
	if path[0] == '.' {
		curPath := utils.GetCurrentDir()
		//fmt.Println("app current Dir: ", curPath)
		pathAfter, find := strings.CutPrefix(path, ".")
		if find {
			path = curPath + pathAfter
		} else {
			path = curPath + "/log/"
		}
	}
	if !utils.DirExists(path) {
		os.Mkdir(path, 0755)
	}

	filename = path + filename
	//fmt.Println("log file name: ", filename)
	return filename
}

func checkIfneedInit() {
	if FileHdl == nil {
		DoInit()
		return
	}

	if utils.FileExists(getLogFileName()) {
		return
	} else {
		FileHdl.Close()
		FileHdl = nil
		DoInit()
	}
}

func getcaller() (file string, line int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file, line
}

func D(v ...any) {
	if beego.BConfig.RunMode == "production" {
		return
	}
	checkIfneedInit()
	file, line := getcaller()
	s := fmt.Sprint(v...)
	DebugLogger.Printf("%s:%d %s\n", file, line, s)

	if beego.BConfig.RunMode == "dev" {
		fmt.Println(v...)
	}
}

func I(v ...any) {
	checkIfneedInit()
	file, line := getcaller()
	s := fmt.Sprint(v...)
	InfoLogger.Printf("%s:%d %s\n", file, line, s)

	if beego.BConfig.RunMode == "dev" {
		fmt.Println(v...)
	}
}

func W(v ...any) {
	checkIfneedInit()
	file, line := getcaller()
	s := fmt.Sprint(v...)
	WarningLogger.Printf("%s:%d %s\n", file, line, s)

	if beego.BConfig.RunMode == "dev" {
		fmt.Println(v...)
	}
}

func E(v ...any) {
	checkIfneedInit()
	file, line := getcaller()
	s := fmt.Sprint(v...)
	ErrorLogger.Printf("%s:%d %s\n", file, line, s)

	if beego.BConfig.RunMode == "dev" {
		fmt.Println(v...)
	}
}
