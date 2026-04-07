package logger

import (
	"errors"

	"github.com/rs/zerolog"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
	FATAL
)

func SetLevel(level LogLevel) error {
	switch level {
	case INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case WARN:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case FATAL:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		return errors.New("unknown log level")
	}
	return nil
}
