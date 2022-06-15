package content

import (
	"fmt"
	"tec-doc/internal/config"
	"tec-doc/internal/model"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

//jRPC клиент для работы с контентом
type ContentClientgRPC interface {
}

type contentClientgRPC struct {
	cfg *config.Config
	log	*zerolog.Logger
}

func (c contentClientgRPC) PushContent(contents []model.Content) error {
	conn, err := grpc.Dial(c.cfg.ContentClientConfig.Url, 
				grpc.WithTimeout(c.cfg.ContentClientConfig.Timeout), 
				grpc.WithInsecure)

	if err != nil {
		return fmt.Errorf("can't dial grpc")
	}

	defer conn.Close()

	return nil
}

//??????