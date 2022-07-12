package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

func InitLogger(level string) (zerolog.Logger, error) {
	res, err := zerolog.ParseLevel(level)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("getting log level error [%w]", err)
	}
	return zerolog.New(os.Stdout).Level(res).With().Timestamp().Logger(), nil
}
