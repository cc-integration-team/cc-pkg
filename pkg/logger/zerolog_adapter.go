package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	ctx context.Context
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
		level, ok := mapZerologLevel[strings.ToLower(cfg.File.Level)]
		if !ok {
			level = zerolog.InfoLevel
		}
		filteredFileWriter := &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: fileWriter},
			Level:  level,
		}
		writers = append(writers, filteredFileWriter)
	}

	// if console is enabled or there is no file writer, add console writer
	if cfg.Console.Enabled || len(writers) == 0 {
		level, ok := mapZerologLevel[strings.ToLower(cfg.Console.Level)]
		if !ok {
			level = zerolog.InfoLevel
		}
		filteredConsoleWriter := &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: os.Stdout},
			Level:  level,
		}
		writers = append(writers, filteredConsoleWriter)
	}

	log := zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().CallerWithSkipFrameCount(4).Logger()

	return &zerologAdapter{log: &log}
}

func (l *zerologAdapter) Debug(msg string) {
	l.log.Debug().Ctx(l.ctx).Msg(msg)
}

func (l *zerologAdapter) Debugf(format string, v ...any) {
	l.log.Debug().Ctx(l.ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Info(msg string) {
	l.log.Info().Ctx(l.ctx).Msg(msg)
}

func (l *zerologAdapter) Infof(format string, v ...any) {
	l.log.Info().Ctx(l.ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Warn(msg string) {
	l.log.Warn().Ctx(l.ctx).Msg(msg)
}

func (l *zerologAdapter) Warnf(format string, v ...any) {
	l.log.Warn().Ctx(l.ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Error(msg string) {
	l.log.Error().Ctx(l.ctx).Msg(msg)
}

func (l *zerologAdapter) Errorf(format string, v ...any) {
	l.log.Error().Ctx(l.ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Fatal(msg string) {
	l.log.Fatal().Ctx(l.ctx).Msg(msg)
}

func (l *zerologAdapter) Fatalf(format string, v ...any) {
	l.log.Fatal().Ctx(l.ctx).Msgf(format, v...)
}

func (l *zerologAdapter) WithContext(ctx context.Context) Logger {
	return &zerologAdapter{log: l.log, ctx: ctx}
}
