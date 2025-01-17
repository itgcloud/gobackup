package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	zerolog.Logger
}

var (
	TimeFormat   = time.RFC3339
	sharedLogger Logger
)

func init() {
	SetLogger("")
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func SetLogger(logPath string) {
	var w io.Writer
	w = zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    TimeFormat,
		PartsOrder:    []string{"time", "level", "tag", "message"},
		FieldsExclude: []string{"tag"},
	}

	if logPath != "" {
		logfile, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		w = zerolog.MultiLevelWriter(logfile, w)
	}

	sharedLogger = newLogger(w)
}

func newLogger(w io.Writer) Logger {
	l := log.Output(w)
	if isDebug() {
		l = l.With().Caller().Logger()
	}
	return Logger{l}
}

func Tag(tag string) Logger {
	return sharedLogger.Tag(fmt.Sprintf("[%s]", tag))
}

func (logger Logger) Tag(tag string) Logger {
	return Logger{logger.With().Str("tag", tag).Logger()}
}

// Debug log
func (logger Logger) Debug(v ...interface{}) {
	logger.Logger.Debug().Msg(sprint(v))
}

// Debugf log
func (logger Logger) Debugf(format string, v ...interface{}) {
	logger.Logger.Debug().Msgf(format, v...)
}

// Info log
func (logger Logger) Info(v ...interface{}) {
	a := sprint(v)
	logger.Logger.Info().Msg(a)
}

// Infof log
func (logger Logger) Infof(format string, v ...interface{}) {
	logger.Logger.Info().Msgf(format, v...)
}

// Warn log
func (logger Logger) Warn(v ...interface{}) {
	logger.Logger.Warn().Msgf(sprint(v))
}

// Warnf log
func (logger Logger) Warnf(format string, v ...interface{}) {
	logger.Logger.Warn().Msgf(format, v...)
}

// Error log
func (logger Logger) Error(v ...interface{}) {
	logger.Logger.Error().Msg(sprint(v))
}

// Errorf log
func (logger Logger) Errorf(format string, v ...interface{}) {
	logger.Logger.Error().Msgf(format, v...)
}

// Fatal log
func (logger Logger) Fatal(v ...interface{}) {
	logger.Logger.Fatal().Msg(sprint(v))
}

// Fatalf log
func (logger Logger) Fatalf(format string, v ...interface{}) {
	logger.Logger.Fatal().Msgf(format, v...)
}

// Debug log
func Debug(v ...interface{}) {
	sharedLogger.Debug(v...)
}

// Info log
func Info(v ...interface{}) {
	sharedLogger.Info(v...)
}

// Warn log
func Warn(v ...interface{}) {
	sharedLogger.Warn(v...)
}

// Error log
func Error(v ...interface{}) {
	sharedLogger.Error(v...)
}

// Errorf log
func Errorf(format string, v ...interface{}) {
	sharedLogger.Errorf(format, v...)
}

// Fatal log
func Fatal(v ...interface{}) {
	sharedLogger.Fatal(v...)
}

func sprint(v ...interface{}) string {
	parts := make([]string, len(v))
	for i, item := range v {
		switch v := item.(type) {
		case []interface{}:
			// Recursively handle nested slices
			parts[i] = sprint(v...)
		default:
			parts[i] = fmt.Sprint(item)
		}
	}
	return strings.Join(parts, " ")
}
