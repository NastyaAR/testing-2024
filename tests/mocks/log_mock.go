package mock_domain

type ILogMock interface {
	Warn(msg string, fields ...any)
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Debug(msg string, fields ...any)
}
