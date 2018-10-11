package stdout

import (
	"os"

	"github.com/yizha/go/logging/writer"
)

// LogWriterCreator is an alias to os.Stdout
type LogWriterCreator struct{}

// New returns a stdout LogWriterCreator
func New() writer.LogWriterCreator {
	return &LogWriterCreator{}
}

// Create returns os.Stdout
func (w *LogWriterCreator) Create(id string) writer.LogWriter {
	return os.Stdout
}
