package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}).
		With().
		Timestamp().
		Logger()
}

func SetFileAndLevel(logFile string, logLevel string) {
	switch logLevel {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Initialize the logger
	if logFile == "" {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.DateTime,
		}).
			With().
			Timestamp().
			Logger()
	} else {
		file, err := os.OpenFile(
			logFile,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		cobra.CheckErr(err)

		var console io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.DateTime,
		}
		output := zerolog.MultiLevelWriter(console, file)

		logger = zerolog.New(output).
			With().
			Timestamp().
			Logger()

		Infof("Logging to: '%s', with log level '%s'", logFile, logLevel)
	}
}

func Print(msg string) {
	logger.Print(msg)
}

func Printf(msg string, v ...interface{}) {
	logger.Printf(msg, v...)
}

func Println(msg string) {
	logger.Println(msg)
}

func Info(msg string) {
	logger.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
}

func Warn(msg string) {
	logger.Warn().Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	logger.Warn().Msgf(format, v...)
}

func Error(msg string, err error) {
	if err == nil {
		logger.Error().Msg(msg)
	} else {
		logger.Error().Err(err).Msg(msg)
	}
}

func Errorf(format string, err error, v ...interface{}) {
	if err == nil {
		logger.Error().Msgf(format, v...)
	} else {
		logger.Error().Err(err).Msgf(format, v...)
	}
}

func Fatal(msg string) {
	logger.Fatal().Msg(msg)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal().Msgf(format, v...)
}

func Debug(msg string) {
	logger.Debug().Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}
