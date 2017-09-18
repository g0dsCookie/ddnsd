package log

import (
	"errors"
	"sync"
)

type Logger interface {
	Close() error

	Emerg(s ...string) error
	Alert(s ...string) error
	Crit(s ...string) error
	Err(s ...string) error
	Warn(s ...string) error
	Notice(s ...string) error
	Info(s ...string) error
	Debug(s ...string) error

	Write(b []byte, s LogSeverity) error
	WriteString(s LogSeverity, m ...string) error
}

var (
	logger    Logger
	loggerMux sync.Mutex
)

func Apply(c LogConfig) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()

	if logger != nil {
		err := logger.Close()
		if err != nil {
			return err
		}
	}

	switch c.Target {
	case LogTargetFile:
		l, err := newFileLogger(c.File, c.Severity)
		if err != nil {
			return err
		}
		logger = l

	case LogTargetSyslog:
		network, addr, err := c.ParseSyslogAddress()
		if err != nil {
			return err
		}
		l, err := newSysLogger(network, addr, c.SyslogTag, c.Facility, c.Severity)
		if err != nil {
			return err
		}
		logger = l

	default:
		return errors.New("invalid log target: " + string(c.Target))
	}

	return nil
}

func Close() error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	if logger != nil {
		return logger.Close()
	}
	return nil
}

func Emerg(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Emerg(s...)
}

func Alert(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Alert(s...)
}

func Crit(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Crit(s...)
}

func Err(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Err(s...)
}

func Warn(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Warn(s...)
}

func Notice(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Notice(s...)
}

func Info(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Info(s...)
}

func Debug(s ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Debug(s...)
}

func Write(b []byte, s LogSeverity) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.Write(b, s)
}

func WriteString(s LogSeverity, m ...string) error {
	loggerMux.Lock()
	defer loggerMux.Unlock()
	return logger.WriteString(s, m...)
}
