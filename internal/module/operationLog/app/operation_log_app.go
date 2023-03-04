package app

import (
	"context"
	pbBasic "github.com/bighuangbee/basic-service/api/basic/v1"
	"github.com/bighuangbee/basic-service/internal/module/operationLog/service"
	"github.com/go-kratos/kratos/v2/log"
)

type OperationLogApp struct {
	pbBasic.UnimplementedOperationLogServer
	uc     *service.OperationLogService
	logger *log.Helper
}

func NewOperationLogApp(uc *service.OperationLogService, logger log.Logger) pbBasic.OperationLogServer {
	return &OperationLogApp{
		uc:     uc,
		logger: log.NewHelper(logger),
	}
}

func (s *OperationLogApp) Add(ctx context.Context, req *pbBasic.AddRequest) (*pbBasic.AddReply, error) {
	s.logger.Error("----------------OperationLogApp", req.Log.UserName)
	return &pbBasic.AddReply{}, nil
}

func (s *OperationLogApp) ListOperationLog(ctx context.Context, req *pbBasic.ListOperationLogRequest) (*pbBasic.ListOperationLogReply, error) {

	return &pbBasic.ListOperationLogReply{}, nil
}

func (s *OperationLogApp) ListOperationLogUser(ctx context.Context, req *pbBasic.ListOperationLogUserRequest) (*pbBasic.ListOperationLogUserReply, error) {

	return &pbBasic.ListOperationLogUserReply{}, nil
}
