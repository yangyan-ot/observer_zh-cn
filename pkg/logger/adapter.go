package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Adapter struct {
	zerolog.Logger
}

func (a *Adapter) Infof(format string, args ...any) {
	a.Logger.Info().Msgf(format, args...)
}

func (a *Adapter) Infoln(args ...any) {
	a.Logger.Info().Msg(fmt.Sprint(args...))
}

func (a *Adapter) Warnf(format string, args ...any) {
	a.Logger.Warn().Msgf(format, args...)
}

func (a *Adapter) Warnln(args ...any) {
	a.Logger.Warn().Msg(fmt.Sprint(args...))
}

func (a *Adapter) Errorf(format string, args ...any) {
	a.Logger.Error().Msgf(format, args...)
}

func (a *Adapter) Errorln(args ...any) {
	a.Logger.Error().Msg(fmt.Sprint(args...))
}

func (a *Adapter) Fatalf(format string, args ...any) {
	a.Logger.Fatal().Msgf(format, args...)
}

func (a *Adapter) Fatalln(args ...any) {
	a.Logger.Fatal().Msg(fmt.Sprint(args...))
}
