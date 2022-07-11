package sig

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var ErrShutdownSignal = errors.New("shutdown")

func Listen(ctx context.Context) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case s := <-shutdown:
		log.Info().Str("signal", s.String()).Msg("signal received")
		return ErrShutdownSignal
	}
}
