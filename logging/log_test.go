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
	CreateLogger("stdout", zerolog.InfoLevel, stdout.New())
	lg := GetLogger("stdout")
	lg.Info().Msg("stdout")
}
