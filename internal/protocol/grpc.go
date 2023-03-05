package protocol

import (
	"time"

	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(bc *conf.Bootstrap, logger log.Logger, services *PbServer) *grpc.Server {
	c := bc.Server

	srv := grpc.NewServer(
		grpc.Address(c.Grpc.Addr),
		grpc.Timeout(time.Duration(c.Grpc.Timeout)*time.Second),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.Validator(),
		),
		grpc.Logger(logger),
	)

	services.RegisterRPC(srv)
	return srv
}
