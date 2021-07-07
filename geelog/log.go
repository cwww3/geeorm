package geelog

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// 字体颜色 https://blog.csdn.net/u014470361/article/details/81512330
	// https://blog.csdn.net/Primeprime/article/details/79708373
	errorLog = log.New(os.Stdout,"\033[4;40;31m[error]\033[0m",log.LstdFlags|log.Lshortfile)
	infoLog = log.New(os.Stdout,"\033[34m[info ]\033[0m",log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu sync.Mutex
)

var (
	Error = errorLog.Println
	Errorf = errorLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
}