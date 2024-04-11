package domain

type Logger interface {
	Info(str string)
	Error(str string)
}
