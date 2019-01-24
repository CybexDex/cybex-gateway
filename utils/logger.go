package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

/*
	日志模块，按日期记录日志
*/
var (
	logPath string
	logFile *os.File
	logger  *Logger
)

const (
	callDepth = 3
)

func init() {
	logger = NewLogger(os.Stderr, "", log.LstdFlags|log.Lshortfile, "DEBUG")
	logger.SetCallDepth(callDepth)
}

// InitLog init logger with log dir and log level
func InitLog(logDir, logLevel string) {
	if len(logLevel) == 0 {
		logLevel = "DEBUG"
	}
	if len(logDir) == 0 {
		logger = NewLogger(os.Stderr, "", log.LstdFlags|log.Lshortfile, logLevel)
		logger.SetCallDepth(callDepth)
		return
	}

	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	logPath = filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		panic(err)
	}
	logger = NewLogger(file, "", log.LstdFlags|log.Lshortfile, logLevel)

	if logger == nil {
		panic("new logger error")
	}
	logger.SetCallDepth(callDepth)

	if err := syscall.Dup2(int(file.Fd()), 1); err != nil {
		panic(err)
	}
	if err := syscall.Dup2(int(file.Fd()), 2); err != nil {
		panic(err)
	}
	logFile = file

	go scheduleDailyRotate(logDir)
}

// GetLogger returns logger reference
func GetLogger() *Logger {
	return logger
}

//call by goroutine
func scheduleDailyRotate(logDir string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("recover in %v, stack: %s\n", r, debug.Stack())
		}
	}()

	for {
		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0,
			now.Nanosecond(), now.Location())
		duration := t.Sub(now)
		time.Sleep(duration)
		switchLogFile(logDir, time.Now())
	}
}

func switchLogFile(logDir string, now time.Time) {
	if logger == nil {
		return
	}
	var err error
	logName := now.Format("2006-01-02") + ".log"
	newFile, err := os.OpenFile(filepath.Join(logDir, logName),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		logger.Println(err)
		return
	}
	syscall.Dup2(int(newFile.Fd()), 1)
	syscall.Dup2(int(newFile.Fd()), 2)
	logger.SetOutput(newFile)

	oldFile := *logFile
	go closeOldFile(&oldFile)
	logFile = newFile
}

func closeOldFile(f *os.File) {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("recover in %v, stack: %s\n", r, debug.Stack())
		}
	}()

	//delay to avoid log lost
	time.Sleep(time.Second * 10)
	err := f.Close()
	if err != nil {
		logger.Println(err.Error())
	}
}

// Debugf print debug log
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

// Infof print info log
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

// Warningf print warning log
func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

// Errorf print error log
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

// Fatalf print fatal log
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

// Debugln print debug log
func Debugln(v ...interface{}) {
	logger.Debugln(v...)
}

// Infoln print info log
func Infoln(v ...interface{}) {
	logger.Infoln(v...)
}

// Warningln print warning log
func Warningln(v ...interface{}) {
	logger.Warningln(v...)
}

// Errorln print error log
func Errorln(v ...interface{}) {
	logger.Errorln(v...)
}

// Fatalln print fatal log
func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

const (
	// DEBUG level
	DEBUG = "DEBUG"
	// INFO level
	INFO = "INFO"
	// WARNING level
	WARNING = "WARNING"
	// ERROR level
	ERROR = "ERROR"
	// FATAL level
	FATAL = "FATAL"
)

const (
	debugFlag   = 1
	infoFlag    = 2
	warningFlag = 3
	errorFlag   = 4
	fatalFlag   = 5
)

var levelMap = map[string]int{
	DEBUG:   debugFlag,
	INFO:    infoFlag,
	WARNING: warningFlag,
	ERROR:   errorFlag,
	FATAL:   fatalFlag,
}

var calldepth = 2

// Logger logger
type Logger struct {
	*log.Logger
	level int
}

// NewLogger return new logger
func NewLogger(out io.Writer, prefix string, flag int, level string) *Logger {
	level = strings.ToUpper(level)
	if _, ok := levelMap[level]; !ok {
		return nil
	}
	return &Logger{log.New(out, prefix, flag), levelMap[level]}
}

// SetOutput set new output
func (l *Logger) SetOutput(out io.Writer) {
	if out == nil {
		return
	}
	l.Logger.SetOutput(out)
}

// SetLevel set new log level
func (l *Logger) SetLevel(level string) {
	if ilvl, ok := levelMap[level]; ok {
		l.level = ilvl
	}
}

// SetCallDepth set new calldepth
func (l *Logger) SetCallDepth(depth int) {
	if depth > 0 {
		calldepth = depth
	}
}

// GetLevel return log level int
func (l *Logger) GetLevel() int {
	return l.level
}

// GetLevelString get log level string
func (l *Logger) GetLevelString() string {
	switch l.level {
	case debugFlag:
		return DEBUG
	case infoFlag:
		return INFO
	case warningFlag:
		return WARNING
	case errorFlag:
		return ERROR
	case fatalFlag:
		return FATAL
	}

	return "WRONG_LEVEL"
}

// Debugf print debug log
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level > debugFlag {
		return
	}
	if len(v) == 0 {
		l.Output(calldepth, fmt.Sprintln("[DEBUG] "+format))
	} else {
		l.Output(calldepth, fmt.Sprintf("[DEBUG] "+format, v...))
	}
}

// Infof print info log
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level > infoFlag {
		return
	}
	if len(v) == 0 {
		l.Output(calldepth, fmt.Sprintln("[INFO] "+format))
	} else {
		l.Output(calldepth, fmt.Sprintf("[INFO] "+format, v...))
	}
}

// Warningf print warning log
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.level > warningFlag {
		return
	}
	if len(v) == 0 {
		l.Output(calldepth, fmt.Sprintln("[WARNING] "+format))
	} else {
		l.Output(calldepth, fmt.Sprintf("[WARNING] "+format, v...))
	}
}

// Errorf print error log
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level > errorFlag {
		return
	}
	if len(v) == 0 {
		l.Output(calldepth, fmt.Sprintln("[ERROR] "+format))
	} else {
		l.Output(calldepth, fmt.Sprintf("[ERROR] "+format, v...))
	}
}

// Fatalf print fatal log
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level > fatalFlag {
		return
	}
	if len(v) == 0 {
		l.Output(calldepth, fmt.Sprintln("[FATAL] "+format))
	} else {
		l.Output(calldepth, fmt.Sprintf("[FATAL] "+format, v...))
	}
}

// Debugln print debug log
func (l *Logger) Debugln(v ...interface{}) {
	if l.level > debugFlag {
		return
	}
	l.Output(calldepth, fmt.Sprintf("%s%s", "[DEBUG] ", fmt.Sprintln(v...)))
}

// Infoln print info log
func (l *Logger) Infoln(v ...interface{}) {
	if l.level > debugFlag {
		return
	}
	l.Output(calldepth, fmt.Sprintf("%s%s", "[Info] ", fmt.Sprintln(v...)))
}

// Warningln print warning log
func (l *Logger) Warningln(v ...interface{}) {
	if l.level > debugFlag {
		return
	}
	l.Output(calldepth, fmt.Sprintf("%s%s", "[WARNING] ", fmt.Sprintln(v...)))
}

// Errorln print error log
func (l *Logger) Errorln(v ...interface{}) {
	if l.level > debugFlag {
		return
	}
	l.Output(calldepth, fmt.Sprintf("%s%s", "[ERROR] ", fmt.Sprintln(v...)))
}

// Fatalln print fatal log
func (l *Logger) Fatalln(v ...interface{}) {
	if l.level > debugFlag {
		return
	}
	l.Output(calldepth, fmt.Sprintf("%s%s", "[FATAL] ", fmt.Sprintln(v...)))
}
