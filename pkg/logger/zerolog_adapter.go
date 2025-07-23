package logger

import (
	"fmt"
	"io"
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
	log *zerolog.Logger
}

func NewZerologAdapter(cfg LoggerConfig) *zerologAdapter {
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

	return &zerologAdapter{log: &log}
}

func (l *zerologAdapter) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *zerologAdapter) Debugf(format string, v ...any) {
	l.log.Debug().Msgf(format, v...)
}

func (l *zerologAdapter) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *zerologAdapter) Infof(format string, v ...any) {
	l.log.Info().Msgf(format, v...)
}

func (l *zerologAdapter) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *zerologAdapter) Warnf(format string, v ...any) {
	l.log.Warn().Msgf(format, v...)
}

func (l *zerologAdapter) Error(msg string) {
	l.log.Error().Msg(msg)
}

func (l *zerologAdapter) Errorf(format string, v ...any) {
	l.log.Error().Msgf(format, v...)
}

func (l *zerologAdapter) Fatal(msg string) {
	l.log.Fatal().Msg(msg)
}

func (l *zerologAdapter) Fatalf(format string, v ...any) {
	l.log.Fatal().Msgf(format, v...)
}
