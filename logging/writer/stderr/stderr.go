package stderr

import (
	"os"

	"github.com/yizha/go/logging/writer"
)

// LogWriterCreator is an alias to os.Stderr
type LogWriterCreator struct{}

// Create returns os.Stdout
func (w *LogWriterCreator) Create(id string) writer.LogWriter {
	return os.Stderr
}
