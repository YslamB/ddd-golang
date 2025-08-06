package logging

import (
	"fmt"
	infra_config "gddd/internal/infrastructure/config"
	"io"
	"log"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	mu       sync.Mutex
	zerolog  zerolog.Logger
	minLevel LogLevel
}

func InitLogger(cfg *infra_config.Config) *Logger {
	var logOutput io.Writer
	if *cfg.IsDebug {
		logOutput = os.Stdout
		fmt.Println("Application running in DEBUG mode. Logs will be printed to terminal.")
	} else {

		if err := os.MkdirAll(cfg.Log.Path, 0755); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}

		logFile, err := os.OpenFile(fmt.Sprintf("%s/%s", cfg.Log.Path, cfg.Log.Filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		logOutput = logFile
		fmt.Printf("Application running in RELEASE mode. Logs will be written to '%s/%s'.\n", cfg.Log.Path, cfg.Log.Filename)
	}

	return NewLogger(logOutput)
}

func NewLogger(output io.Writer) *Logger {

	zl := zerolog.New(output).With().Timestamp().Logger()

	return &Logger{
		zerolog:  zl,
		minLevel: LevelInfo,
	}
}

func (l *Logger) SetMinLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level

}

func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}

	switch level {
	case LevelDebug:
		l.zerolog.Debug().Msgf(format, v...)
	case LevelInfo:
		l.zerolog.Info().Msgf(format, v...)
	case LevelWarn:
		l.zerolog.Warn().Msgf(format, v...)
	case LevelError:
		l.zerolog.Error().Msgf(format, v...)
	case LevelFatal:
		l.zerolog.Fatal().Msgf(format, v...)
	default:
		l.zerolog.Info().Msgf(format, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(LevelDebug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(LevelInfo, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(LevelWarn, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(LevelError, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logf(LevelFatal, format, v...)
}
