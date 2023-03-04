package protocol

import (
	pbAccount "github.com/bighuangbee/basic-service/api/account/v1"
	pbBasic "github.com/bighuangbee/basic-service/api/basic/v1"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type PbServer struct {
	Account pbAccount.AccountServer
	OpLog pbBasic.OperationLogServer
}

func (s *PbServer) RegisterHTTP(srv *http.Server) {
	pbAccount.RegisterAccountHTTPServer(srv, s.Account)
	pbBasic.RegisterOperationLogHTTPServer(srv, s.OpLog)
}

func (s *PbServer) RegisterRPC(srv *grpc.Server) {
	pbAccount.RegisterAccountServer(srv, s.Account)
	pbBasic.RegisterOperationLogServer(srv, s.OpLog)
}
