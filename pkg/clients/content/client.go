package content

import (
	"tec-doc/internal/config"

	// "github.com/gogo/protobuf/protoc-gen-gogo/grpc"
	"github.com/rs/zerolog"
)

//jRPC клиент для работы с контентом

type ContentClientgRPC interface {
}

type contentClientgRPC struct {

	cfg *config.Config
	log	*zerolog.Logger
}

func NewConn()
