package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(path, filename string) *zerolog.Logger {
	file, err := os.OpenFile(path+"/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	var logger zerolog.Logger

	if err != nil {
		// Fallback to stderr if file can't be opened
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		logger.Warn().Err(err).Msg("Failed to log to file, using default stderr")
	} else {
		logger = zerolog.New(file).With().Timestamp().Logger()
	}
	return &logger
}
