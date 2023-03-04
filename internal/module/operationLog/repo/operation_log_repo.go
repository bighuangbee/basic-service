package repo

import (
	"context"
	"github.com/bighuangbee/basic-service/internal/data"
	"github.com/bighuangbee/basic-service/internal/domain"
)

type OperationLogRepo struct {
	data *data.Data
}

func NewOperationLogRepo(data *data.Data) domain.IOperationLogRepo {
	return &OperationLogRepo{
		data: data,
	}
}

func (r *OperationLogRepo) ListOperationLogUser(ctx context.Context, userName string) (user []*domain.UserInfo, err error) {
	return nil, err
}

func (r *OperationLogRepo) Add(ctx context.Context, oplog *domain.OperationLog) error {
	return nil
}

func (r *OperationLogRepo) ListOperationLog(ctx context.Context, query *domain.ListOperationLogRequest) ([]*domain.OperationLog, int32, error) {
	return nil, 0, nil
}
