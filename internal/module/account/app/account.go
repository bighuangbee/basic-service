package app

import (
	"context"
	pb "github.com/bighuangbee/basic-service/api/account/v1"
	"github.com/bighuangbee/basic-service/internal/module/account/service"
	"github.com/go-kratos/kratos/v2/log"
)

type AccountApp struct {
	pb.UnimplementedAccountServer
	svc    *service.AccountService
	logHelper *log.Helper
}

func NewAccountApp(svc *service.AccountService, logHelper *log.Helper) pb.AccountServer {

	return &AccountApp{
		svc:    svc,
		logHelper: logHelper,
	}
}

func (s *AccountApp) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	return s.svc.Login(ctx, req)
}

func (s *AccountApp) Test(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	s.logHelper.Debug("112233")
	return &pb.LoginReply{UserId: 10086}, nil
}
