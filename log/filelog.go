package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LoggerNotAvailable struct{}

func (e LoggerNotAvailable) Error() string { return "logger is not available" }

type fileLogger struct {
	file      *os.File
	logstream *log.Logger
	severity  LogSeverity
}

func newFileLogger(file string, severity LogSeverity) (f *fileLogger, err error) {
	fileHandle, err := os.Create(file)
	if err != nil {
		return
	}

	f = &fileLogger{
		file:      fileHandle,
		logstream: log.New(fileHandle, "", log.LstdFlags),
		severity:  severity,
	}
	return
}

func (f *fileLogger) Close() error {
	if f.file == nil {
		return nil
	}
	f.logstream = nil
	err := f.file.Close()
	f.file = nil
	return err
}

func (f *fileLogger) WriteString(s LogSeverity, m ...string) error {
	if f.logstream == nil {
		return LoggerNotAvailable{}
	}
	if f.severity.Bigger(s) {
		return nil
	}
	prefix := fmt.Sprintf("[%v]", strings.ToUpper(string(s)))
	for _, v := range m {
		f.logstream.Printf("%v %v\n", prefix, v)
	}
	return nil
}

func (f *fileLogger) Write(b []byte, s LogSeverity) error { return f.WriteString(s, string(b)) }
func (f *fileLogger) Emerg(s ...string) error             { return f.WriteString(LogSeverityEmerg, s...) }
func (f *fileLogger) Alert(s ...string) error             { return f.WriteString(LogSeverityAlert, s...) }
func (f *fileLogger) Crit(s ...string) error              { return f.WriteString(LogSeverityCrit, s...) }
func (f *fileLogger) Err(s ...string) error               { return f.WriteString(LogSeverityErr, s...) }
func (f *fileLogger) Warn(s ...string) error              { return f.WriteString(LogSeverityWarn, s...) }
func (f *fileLogger) Notice(s ...string) error            { return f.WriteString(LogSeverityNotice, s...) }
func (f *fileLogger) Info(s ...string) error              { return f.WriteString(LogSeverityInfo, s...) }
func (f *fileLogger) Debug(s ...string) error             { return f.WriteString(LogSeverityDebug, s...) }
