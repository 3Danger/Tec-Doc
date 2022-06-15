package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"tec-doc/internal/config"
	"time"
)

func NewLogger(conf *config.Config) zerolog.Logger {
	// Define output format
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s |", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return "\033[36m" + fmt.Sprintf("%s", i) + "\033[0m"
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	// New logger
	log := zerolog.New(output).With().Timestamp().Logger()

	// Setting log level
	switch strings.ToLower(conf.LogLevel) {
	case "debug", "d":
		log.Level(zerolog.DebugLevel)
	case "info", "i":
		log.Level(zerolog.InfoLevel)
	case "fatal", "f":
		log.Level(zerolog.FatalLevel)
	case "error", "err", "e":
		log.Level(zerolog.ErrorLevel)
	case "panic", "pan", "p":
		log.Level(zerolog.PanicLevel)
	case "trace", "t":
		log.Level(zerolog.TraceLevel)
	case "warning", "warn", "w":
		log.Level(zerolog.WarnLevel)
	}
	return log
}
