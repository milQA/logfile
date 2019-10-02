package logfile

import (
	"log"
	"os"
	"sync"
	"time"
)

type LogSaver interface {
	Write(b []byte) (n int, err error)
	Close()
}
type LogSave struct {
	mu          sync.Mutex
	file        *os.File
	duplicate   bool
	logfileName string
}

// need defer *LogSave.Close()
func NewLogSave(logfileName string, logSaveTick time.Duration) *LogSave {
	timePrefix := time.Now().Format("2006-01-02_15-04")
	fileName := logfileName + "-" + timePrefix
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logSave := LogSave{
		file:        f,
		duplicate:   true,
		logfileName: logfileName,
	}
	if int64(logSaveTick) != int64(0) {
		go logSave.runDuplicater(logSaveTick)
	}
	return &logSave
}

func (l *LogSave) runDuplicater(logSaveTick time.Duration) {
	for l.duplicate {
		<-time.After(logSaveTick)
		log.Printf("[LogSave] Duplicate logfile")
		l.duplicateFile()
	}
}

func (l *LogSave) duplicateFile() {
	l.mu.Lock()
	defer l.mu.Unlock()
	err := l.file.Close()
	if err != nil {
		log.Fatalf("error closing file: %v", err)
	}
	timePrefix := time.Now().Format("2006-01-02_15-04")
	fileName := l.logfileName + "-" + timePrefix
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	l.file = f
}

func (l *LogSave) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Write(b)
}

func (l *LogSave) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.file.Close()
	l.duplicate = false
}
