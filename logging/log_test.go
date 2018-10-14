package logging

import (
	"testing"

	"github.com/rs/zerolog"
	//"github.com/yizha/go/logging/writer/lumberjack"
	"github.com/yizha/go/logging/writer/stdout"
)

/*func TestStdout(t *testing.T) {
	SetupDefaults(zerolog.InfoLevel, &stdout.LogWriterCreator{})
	lg := GetLogger("stdout")
	lg.Info().Msg("info message")
}*/

func TestLumberjack(t *testing.T) {
	/*SetupDefaults(zerolog.InfoLevel, lumberjack.New("/tmp", 1, 0, 0, false))
	lg := GetLogger("output-1")
	lg.Info().Msg("output-1")
	lg = GetLogger("output-2")
	lg.Info().Msg("output-2")*/
	SetupGlobalConf(DefaultGlobalConf().SetTimestampFormat("2006-01-02T15:04:05.999999"))
	CreateLogger("stdout", zerolog.InfoLevel, true, false, stdout.New())
	lg := GetLogger("stdout")
	lg.Info().Msg("stdout")
}
