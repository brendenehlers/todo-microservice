package todo

type Logger interface {
	Info(strs ...string)
	Error(strs ...string)
}
