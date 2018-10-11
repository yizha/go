package writer

import (
	"io"
)

// A LogWriter is used by zerolog instance to write log data to.
type LogWriter io.Writer

// A LogWriterCreator creates LogWriter.
type LogWriterCreator interface {
	Create(string) LogWriter
}
