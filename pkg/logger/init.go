package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logWriters = []io.Writer{}

func Init() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	logWriters = []io.Writer{}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05.000",
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			moduleKey,
			zerolog.MessageFieldName,
		},
		FieldsExclude: []string{moduleKey},
		FormatPartValueByName: func(i any, name string) string {
			if name == moduleKey {
				return fmt.Sprintf("\033[35m[module:%s]\033[0m", i)
			}
			return fmt.Sprintf("%v", i)
		},
	}
	logWriters = append(logWriters, output)
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(logWriters...)).
		With().
		Timestamp().
		Logger()
}
