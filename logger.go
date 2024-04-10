package todo

type Logger interface {
	Info(str string)
	Error(str string)
}
