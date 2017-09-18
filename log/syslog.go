package log

// +build !windows,!nacl,!plan9

import (
	"errors"
	"log/syslog"
)

const LogTargetSyslog LogTarget = "syslog"

type sysLogger struct {
	logstream *syslog.Writer
}

func newSysLogger(network, addr, tag string, prio LogFacility, severity LogSeverity) (l *sysLogger, err error) {
	priority := prio.Priority() | severity.Priority()
	stream, err := syslog.Dial(network, addr, priority, tag)
	if err != nil {
		return
	}
	l = &sysLogger{logstream: stream}
	return
}

func (l *sysLogger) Close() error {
	if l.logstream == nil {
		return nil
	}
	err := l.logstream.Close()
	l.logstream = nil
	return err
}

func (l *sysLogger) WriteString(s LogSeverity, m ...string) error {
	if l.logstream == nil {
		return LoggerNotAvailable{}
	}
	for _, v := range m {
		var err error
		switch s {
		case LogSeverityEmerg:
			err = l.logstream.Emerg(v)
		case LogSeverityAlert:
			err = l.logstream.Alert(v)
		case LogSeverityCrit:
			err = l.logstream.Crit(v)
		case LogSeverityErr:
			err = l.logstream.Err(v)
		case LogSeverityWarn:
			err = l.logstream.Warning(v)
		case LogSeverityNotice:
			err = l.logstream.Notice(v)
		case LogSeverityInfo:
			err = l.logstream.Info(v)
		case LogSeverityDebug:
			err = l.logstream.Debug(v)
		default:
			err = errors.New("unknown severity " + string(s))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *sysLogger) Write(b []byte, s LogSeverity) error { return l.WriteString(s, string(b)) }
func (l *sysLogger) Emerg(s ...string) error             { return l.WriteString(LogSeverityEmerg, s...) }
func (l *sysLogger) Alert(s ...string) error             { return l.WriteString(LogSeverityAlert, s...) }
func (l *sysLogger) Crit(s ...string) error              { return l.WriteString(LogSeverityCrit, s...) }
func (l *sysLogger) Err(s ...string) error               { return l.WriteString(LogSeverityErr, s...) }
func (l *sysLogger) Warn(s ...string) error              { return l.WriteString(LogSeverityWarn, s...) }
func (l *sysLogger) Notice(s ...string) error            { return l.WriteString(LogSeverityNotice, s...) }
func (l *sysLogger) Info(s ...string) error              { return l.WriteString(LogSeverityInfo, s...) }
func (l *sysLogger) Debug(s ...string) error             { return l.WriteString(LogSeverityDebug, s...) }
