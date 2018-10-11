package lumberjack

import (
	"fmt"
	"os"
	"path"

	"github.com/yizha/go/logging/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogWriterCreator is a writer.LogWriterCreator, it can create
// a lumberjack instance.
type LogWriterCreator struct {
	logDirPath string
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
}

func validateLumberjackConfs(dirPath string, maxSize, maxBackups, maxAge int, compress bool) {
	if dirPath == "" {
		panic("blank log dir path!")
	}
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		panic(fmt.Sprintf("failed to mkdir dirs %s, error: %v", dirPath, err))
	}
	if maxSize < 1 {
		panic(fmt.Sprintf("maxSize %v is less than 1!", maxSize))
	}
	if maxBackups < 0 {
		panic("maxBackups %v is negative!")
	}
	if maxAge < 0 {
		panic("maxAge %v is negative!")
	}
}

// New creates a Lumberjack log writer creator
func New(dirPath string, maxSize, maxBackups, maxAge int, compress bool) *LogWriterCreator {
	validateLumberjackConfs(dirPath, maxSize, maxBackups, maxAge, compress)
	return &LogWriterCreator{
		logDirPath: dirPath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
	}
}

// Create returns an instance of Lumberjack{}
func (w *LogWriterCreator) Create(id string) writer.LogWriter {
	fpath := path.Join(w.logDirPath, fmt.Sprintf("%s.log", id))
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(fmt.Sprintf("failed to open/create %s, error: %v", fpath, err))
	}
	if err = f.Close(); err != nil {
		panic(fmt.Sprintf("failed to close file %s, error: %v", fpath, err))
	}
	return &lumberjack.Logger{
		Filename:   fpath,
		MaxSize:    w.maxSize,
		MaxBackups: w.maxBackups,
		MaxAge:     w.maxAge,
		Compress:   w.compress,
	}
}
