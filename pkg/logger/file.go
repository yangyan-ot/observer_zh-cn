package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func RegisterFileLogger(filePath string, maxSize, rotation, lifeCycle int) {
	logWriters = append(logWriters, &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: rotation,
		MaxAge:     lifeCycle,
		Compress:   true,
	})
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(logWriters...)).With().Timestamp().Logger()
}
