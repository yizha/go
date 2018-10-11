package logging // import "github.com/yizha/go/logging"

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/yizha/go/logging/writer"
	"github.com/yizha/go/logging/writer/stdout"
)

// GlobalField is the type for global field
type GlobalField string

// Global Fields
const (
	TimestampField GlobalField = "field-timestamp"
	LevelField                 = "field-level"
	MessageField               = "field-message"
	ErrorField                 = "field-error"
	CallerField                = "field-caller"
)

// GlobalConf contains global configurations for zerolog which can only
// be changed before any zerolog instance is created.
type GlobalConf struct {
	level                zerolog.Level
	fieldNames           map[GlobalField]string
	timestampFormat      string
	timestampFunc        func() time.Time
	callerSkipFrameCount int
}

// DefaultGlobalFieldNames returns a map containing defaults names
// for global fields:
//   timestamp: "log-ts"
//   level:     "log-level"
//   message:   "message"
//   error:     "error"
//   caller:    "caller"
func DefaultGlobalFieldNames() map[GlobalField]string {
	m := make(map[GlobalField]string)
	m[TimestampField] = "log-ts"
	m[LevelField] = "log-level"
	m[MessageField] = "message"
	m[ErrorField] = "error"
	m[CallerField] = "caller"
	return m
}

// SetLevel sets the global log level for zerolog.
func (conf *GlobalConf) SetLevel(lvl zerolog.Level) *GlobalConf {
	conf.level = lvl
	return conf
}

// SetFieldNames change the name(s) of globl fields.
func (conf *GlobalConf) SetFieldNames(names map[GlobalField]string) *GlobalConf {
	if conf == nil {
		return nil
	}
	if conf.fieldNames == nil {
		conf.fieldNames = DefaultGlobalFieldNames()
	}
	for field, name := range names {
		conf.fieldNames[field] = name
	}
	m := make(map[string]bool)
	for _, name := range conf.fieldNames {
		m[name] = true
	}
	if len(m) < 5 {
		panic(fmt.Sprintf("there are dup field names after setting new names: %v", conf.fieldNames))
	}
	return conf
}

// SetTimestampFormat sets the format string for the timestamp field value.
// The default format is "2006-01-02T15:04:05.999".
func (conf *GlobalConf) SetTimestampFormat(fmt string) *GlobalConf {
	if fmt == "" {
		panic("blank timestamp format!")
	}
	conf.timestampFormat = fmt
	return conf
}

// SetTimestampFunc sets the function which generates the value for
// the timestamp field. Default function returns time.Now().UTC().
func (conf *GlobalConf) SetTimestampFunc(f func() time.Time) *GlobalConf {
	if f == nil {
		panic("nil timestamp function!")
	}
	conf.timestampFunc = f
	return conf
}

// SetCallerSkipFrameCount sets the skip count when getting caller info.
// Default value is 2.
func (conf *GlobalConf) SetCallerSkipFrameCount(n int) *GlobalConf {
	if n < 0 {
		panic(fmt.Sprintf("negative caller skip frame count: %v", n))
	}
	conf.callerSkipFrameCount = n
	return conf
}

// DefaultGlobalConf returns a GlobalConf instance with default values.
func DefaultGlobalConf() *GlobalConf {
	return &GlobalConf{
		level:                zerolog.InfoLevel,
		fieldNames:           DefaultGlobalFieldNames(),
		timestampFormat:      "2006-01-02T15:04:05.999",
		timestampFunc:        func() time.Time { return time.Now().UTC() },
		callerSkipFrameCount: 2,
	}
}

// SetupGlobalConf sets the global configurations of zerolog,
// it can be called multiple times and the last call wins but it
// MUST be called before GetLogger() is called.
func SetupGlobalConf(conf *GlobalConf) {
	lock.Lock()
	defer lock.Unlock()
	if getLoggerCalled {
		lg := GetLogger("warn")
		lg.Info().Msg("ignore call to SetupGlobalConf() as GetLogger() has been invoked.")
	} else {
		zerolog.TimestampFieldName = conf.fieldNames[TimestampField]
		zerolog.LevelFieldName = conf.fieldNames[LevelField]
		zerolog.MessageFieldName = conf.fieldNames[MessageField]
		zerolog.ErrorFieldName = conf.fieldNames[ErrorField]
		zerolog.CallerFieldName = conf.fieldNames[CallerField]
		zerolog.TimeFieldFormat = conf.timestampFormat
		zerolog.TimestampFunc = conf.timestampFunc
		zerolog.CallerSkipFrameCount = conf.callerSkipFrameCount
		zerolog.SetGlobalLevel(conf.level)
	}
}

func init() {
	SetupGlobalConf(DefaultGlobalConf())
}

type loggerDefaults struct {
	level         zerolog.Level
	writerCreator writer.LogWriterCreator
}

var (
	getLoggerCalled = false
	loggers         = make(map[string]*zerolog.Logger)
	lock            = &sync.Mutex{}
	loggerDef       = &loggerDefaults{zerolog.InfoLevel, &stdout.LogWriterCreator{}}
)

// SetupDefaults sets up default settings for logger. Inside GetLogger(), it
// creates logger from this defaults if the requested logger doesn't exist.
func SetupDefaults(lvl zerolog.Level, writerCreator writer.LogWriterCreator) {
	lock.Lock()
	defer lock.Unlock()
	if getLoggerCalled {
		lg := GetLogger("warn")
		lg.Info().Msg("ignore call to setupDefaults() as GetLogger() has been invoked.")
	} else {
		loggerDef = &loggerDefaults{lvl, writerCreator}
	}
}

// CreateLogger creates a logger and stores it for later retrieving.
// This SHOULD be called at app startup to setup loggers.
func CreateLogger(id string, lvl zerolog.Level, wcreator writer.LogWriterCreator) {
	lock.Lock()
	defer lock.Unlock()
	w := wcreator.Create(id)
	lg := zerolog.New(w).Level(lvl).With().Timestamp().Caller().Logger()
	loggers[id] = &lg
}

// GetLogger retrieves an existing logger or creates a new one with
// logger default settings.
func GetLogger(id string) *zerolog.Logger {
	lock.Lock()
	getLoggerCalled = true
	defer lock.Unlock()
	lg, ok := loggers[id]
	if !ok {
		w := loggerDef.writerCreator.Create(id)
		newLogger := zerolog.New(w).Level(loggerDef.level).With().Timestamp().Caller().Logger()
		lg = &newLogger
		loggers[id] = lg
	}
	return lg
}
