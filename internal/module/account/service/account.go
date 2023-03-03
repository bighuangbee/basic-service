package service

import (
	"context"
	pb "github.com/bighuangbee/basic-service/api/account/v1"
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/basic-service/internal/domain"
	basicPb "github.com/bighuangbee/gokit/api/basic/v1"
	kitKratos "github.com/bighuangbee/gokit/kratos"
	"github.com/go-kratos/kratos/v2/log"
)

type AccountService struct {
	repo   domain.IAccountRepo
	logger *log.Helper
	bc     *conf.Bootstrap
}

func NewAccountService(repo domain.IAccountRepo, logger log.Logger, bc *conf.Bootstrap) *AccountService {
	return &AccountService{
		repo:   nil,
		logger: log.NewHelper(logger),
		bc:     nil,
	}
}

func (this *AccountService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {

	if req.Username == "" {
		//return nil, errors.New("123c")
		return nil, kitKratos.ResponseErr(ctx, basicPb.ErrorInvalidParameter)
	}

	return nil, kitKratos.ResponseErr(ctx, pb.ErrorAccountPwdError)

	return &pb.LoginReply{
		UserId: 10088,
	}, nil
}
