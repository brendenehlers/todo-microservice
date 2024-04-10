package slogger

import "log/slog"

type Slogger struct{}

func (*Slogger) Info(str string) {
	slog.Info(str)
}

func (*Slogger) Error(str string) {
	slog.Error(str)
}
