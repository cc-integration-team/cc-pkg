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
			cw := zerolog.ConsoleWriter{
				Out:           fileWriter,
				TimeFormat:    time.RFC3339,
				NoColor:       true,
				FieldsExclude: []string{"service"},
				FormatLevel: func(i any) string {
					return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
				},
			}
			if cfg.Service != "" {
				cw.FormatMessage = func(i any) string {
					return fmt.Sprintf("%s | %v", cfg.Service, i)
				}
			}
			writer = cw
		} else {
			writer = fileWriter
		}

		level, ok := mapZerologLevel[strings.ToLower(cfg.File.Level)]
		if !ok {
			level = zerolog.InfoLevel
		}
		var fileLevelWriter zerolog.LevelWriter = &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: writer},
			Level:  level,
		}
		if cfg.Masking.Enabled {
			fileLevelWriter = newMaskingLevelWriter(fileLevelWriter, cfg.Masking)
		}

		writers = append(writers, fileLevelWriter)
	}

	// if console is enabled or there is no file writer, add console writer
	if cfg.Console.Enabled || len(writers) == 0 {
		var writer io.Writer
		if cfg.Console.Pretty {
			cw := zerolog.ConsoleWriter{
				Out:           os.Stdout,
				TimeFormat:    time.RFC3339,
				NoColor:       true,
				FieldsExclude: []string{"service"},
				FormatLevel: func(i any) string {
					return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
				},
			}
			if cfg.Service != "" {
				cw.FormatMessage = func(i any) string {
					return fmt.Sprintf("%s | %v", cfg.Service, i)
				}
			}
			writer = cw
		} else {
			writer = os.Stdout
		}

		level, ok := mapZerologLevel[strings.ToLower(cfg.Console.Level)]
		if !ok {
			level = zerolog.InfoLevel
		}

		var consoleLevelWriter zerolog.LevelWriter = &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: writer},
			Level:  level,
		}
		if cfg.Masking.Enabled {
			consoleLevelWriter = newMaskingLevelWriter(consoleLevelWriter, cfg.Masking)
		}
		writers = append(writers, consoleLevelWriter)
	}

	var log zerolog.Logger
	ctx := zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp()
	if cfg.Service != "" {
		ctx = ctx.Str("service", cfg.Service)
	}
	if cfg.Caller {
		log = ctx.CallerWithSkipFrameCount(3).Logger()
	} else {
		log = ctx.Logger()
	}

	return &zerologAdapter{log: &log}
}

type zerologAdapter struct {
	log *zerolog.Logger
}

func (l *zerologAdapter) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *zerologAdapter) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *zerologAdapter) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *zerologAdapter) Error(msg string) {
	l.log.Error().Msg(msg)
}

func (l *zerologAdapter) Fatal(msg string) {
	l.log.Fatal().Msg(msg)
}

func (l *zerologAdapter) Debugf(msg string, args ...any) {
	l.log.Debug().Msgf(msg, args...)
}

// Errorf implements [Logger].
func (l *zerologAdapter) Errorf(msg string, args ...any) {
	l.log.Error().Msgf(msg, args...)
}

// Fatalf implements [Logger].
func (l *zerologAdapter) Fatalf(msg string, args ...any) {
	l.log.Fatal().Msgf(msg, args...)
}

// Infof implements [Logger].
func (l *zerologAdapter) Infof(msg string, args ...any) {
	l.log.Info().Msgf(msg, args...)
}

// Warnf implements [Logger].
func (l *zerologAdapter) Warnf(msg string, args ...any) {
	l.log.Warn().Msgf(msg, args...)
}

func (l *zerologAdapter) WithFields(fields Fields) Logger {
	ctx := l.log.With()
	for k, v := range fields {
		if err, ok := v.(error); ok {
			ctx = ctx.AnErr(k, err)
		} else {
			ctx = ctx.Any(k, v)
		}
	}
	newLog := ctx.Logger()
	return &zerologAdapter{log: &newLog}
}
