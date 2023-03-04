package service

import (
	"context"
	"github.com/bighuangbee/basic-service/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

type OperationLogService struct {
	repo   domain.IOperationLogRepo
	logger *log.Helper
}

func NewOperationLogService(repo domain.IOperationLogRepo, logger log.Logger) *OperationLogService {
	return &OperationLogService{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

func (s *OperationLogService) Add(ctx context.Context, oplog *domain.OperationLog) error {
	return s.repo.Add(ctx, oplog)
}

func (s *OperationLogService) ListOperationLog(ctx context.Context, query *domain.ListOperationLogRequest) (
	logs []*domain.OperationLog, total int32, err error) {

	return s.repo.ListOperationLog(ctx, query)
}

func (s *OperationLogService) ListOperationLogUser(ctx context.Context, userName string) (users []*domain.UserInfo, err error) {
	return s.repo.ListOperationLogUser(ctx, userName)
}
