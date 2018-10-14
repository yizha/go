package lumberjack

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yizha/go/logging/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogWriterCreator is a writer.LogWriterCreator, which creats
// a lumberjack instance as writer.LogWriter
type LogWriterCreator struct {
	logDirPath string
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
}

func validateLumberjackConfs(dirPath string, maxSize, maxBackups, maxAge int, compress bool) error {
	if dirPath == "" {
		return fmt.Errorf("blank log dir path")
	}
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to mkdir dirs %s, error: %v", dirPath, err)
	}
	if maxSize < 1 {
		return fmt.Errorf("maxSize %v is less than 1", maxSize)
	}
	if maxBackups < 0 {
		return fmt.Errorf("maxBackups %v is negative", maxBackups)
	}
	if maxAge < 0 {
		return fmt.Errorf("maxAge %v is negative", maxAge)
	}
	return nil
}

// New creates a Lumberjack log writer creator
func New(dirPath string, maxSize, maxBackups, maxAge int, compress bool) *LogWriterCreator {
	if err := validateLumberjackConfs(dirPath, maxSize, maxBackups, maxAge, compress); err != nil {
		panic(err.Error())
	}
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
	fpath := filepath.Join(w.logDirPath, filepath.FromSlash(fmt.Sprintf("%s.log", id)))
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
