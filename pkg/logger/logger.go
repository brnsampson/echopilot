package logger

// This interface can and will probably need to be modified if you don't want to use the uber zap logger.
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Sync() error
}
