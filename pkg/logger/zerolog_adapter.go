package logger

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var mapZerologLevel = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

type zerologAdapter struct {
	log    *zerolog.Logger
	fields Fields
}

func NewZerologAdapter(cfg LoggerConfig) Logger {
	var writers []io.Writer

	// if file is enabled, add file writer
	if cfg.File.Enabled {
		_ = os.MkdirAll(filepath.Dir(cfg.File.Path), 0755)
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.File.Path,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			MaxAge:     cfg.File.MaxAge,
			Compress:   cfg.File.Compress,
		}
		var writer io.Writer
		if cfg.File.Pretty {
			writer = zerolog.ConsoleWriter{
				Out:        fileWriter,
				TimeFormat: time.RFC3339,
				NoColor:    true,
				FormatLevel: func(i any) string {
					return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
				},
			}
		} else {
			writer = fileWriter
		}

		level, ok := mapZerologLevel[strings.ToLower(cfg.File.Level)]
		if !ok {
			level = zerolog.InfoLevel
		}
		filteredFileWriter := &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: writer},
			Level:  level,
		}

		writers = append(writers, filteredFileWriter)
	}

	// if console is enabled or there is no file writer, add console writer
	if cfg.Console.Enabled || len(writers) == 0 {
		var writer io.Writer
		if cfg.Console.Pretty {
			writer = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
				NoColor:    true,
				FormatLevel: func(i any) string {
					return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
				},
			}
		} else {
			writer = os.Stdout
		}

		level, ok := mapZerologLevel[strings.ToLower(cfg.Console.Level)]
		if !ok {
			level = zerolog.InfoLevel
		}

		filteredConsoleWriter := &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: writer},
			Level:  level,
		}
		writers = append(writers, filteredConsoleWriter)
	}

	var log zerolog.Logger
	if cfg.CallerDebug {
		log = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().CallerWithSkipFrameCount(4).Logger()
	} else {
		log = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()
	}

	return &zerologAdapter{log: &log, fields: make(Fields)}
}

func (l *zerologAdapter) Debug(msg string) {
	l.logMsg(l.log.Debug(), msg)
}

func (l *zerologAdapter) Debugf(format string, v ...any) {
	l.logMsgf(l.log.Debug(), format, v...)
}

func (l *zerologAdapter) Info(msg string) {
	l.logMsg(l.log.Info(), msg)
}

func (l *zerologAdapter) Infof(format string, v ...any) {
	l.logMsgf(l.log.Info(), format, v...)
}

func (l *zerologAdapter) Warn(msg string) {
	l.logMsg(l.log.Warn(), msg)
}

func (l *zerologAdapter) Warnf(format string, v ...any) {
	l.logMsgf(l.log.Warn(), format, v...)
}

func (l *zerologAdapter) Error(msg string) {
	l.logMsg(l.log.Error(), msg)
}

func (l *zerologAdapter) Errorf(format string, v ...any) {
	l.logMsgf(l.log.Error(), format, v...)
}

func (l *zerologAdapter) Fatal(msg string) {
	l.logMsg(l.log.Fatal(), msg)
}

func (l *zerologAdapter) Fatalf(format string, v ...any) {
	l.logMsgf(l.log.Fatal(), format, v...)
}

func (l *zerologAdapter) WithFields(fields Fields) Logger {
	newFields := make(Fields)
	maps.Copy(newFields, l.fields)
	maps.Copy(newFields, fields)

	return &zerologAdapter{
		log:    l.log,
		fields: newFields,
	}
}

func (l *zerologAdapter) logMsg(e *zerolog.Event, msg string) {
	for k, v := range l.fields {
		e = e.Any(k, v)
	}
	e.Msg(msg)
}

func (l *zerologAdapter) logMsgf(e *zerolog.Event, format string, v ...any) {
	for k, v := range l.fields {
		e = e.Any(k, v)
	}
	e.Msgf(format, v...)
}
