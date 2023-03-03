package repo

import (
	"context"
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/basic-service/internal/data"
	"github.com/bighuangbee/basic-service/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
)


func NewAccountRepo(data *data.Data, logHelper *log.Helper, bootstrap *conf.Bootstrap) domain.IAccountRepo {
	return &AccountRepo{
		data:   data,
		logHelper: logHelper,
		bc:     bootstrap,
	}
}

type AccountRepo struct {
	data   *data.Data
	logHelper *log.Helper
	bc     *conf.Bootstrap
}

func (this *AccountRepo) Login(context.Context, *domain.Account)  (*domain.Account, error) {
	return &domain.Account{}, nil
}
