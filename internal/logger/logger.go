package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func InitLogger() {
	// Use console writer for better development experience
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	Logger = zerolog.New(output).With().Timestamp().Logger()
	log.Logger = Logger

	// Set global log level (can be overridden by environment variable)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// In development, use debug level
	if os.Getenv("ENV") == "development" || os.Getenv("ENV") == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// Info logs an info message
func Info(msg string) {
	Logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	Logger.Info().Msgf(format, v...)
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	Logger.Debug().Msgf(format, v...)
}

// Error logs an error message
func Error(msg string) {
	Logger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	Logger.Error().Msgf(format, v...)
}

// ErrorErr logs an error with error details
func ErrorErr(err error, msg string) {
	Logger.Error().Err(err).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func Warnf(format string, v ...interface{}) {
	Logger.Warn().Msgf(format, v...)
}

// Fatal logs a fatal error and exits
func Fatal(msg string) {
	Logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal error and exits
func Fatalf(format string, v ...interface{}) {
	Logger.Fatal().Msgf(format, v...)
}

// WithField returns a logger with a field
func WithField(key string, value interface{}) *zerolog.Event {
	return Logger.Info().Interface(key, value)
}

// WithError returns a logger with an error field
func WithError(err error) *zerolog.Event {
	return Logger.Error().Err(err)
}
