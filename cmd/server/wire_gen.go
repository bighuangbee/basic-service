// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/basic-service/internal/data"
	"github.com/bighuangbee/basic-service/internal/module/account/app"
	"github.com/bighuangbee/basic-service/internal/module/account/repo"
	"github.com/bighuangbee/basic-service/internal/module/account/service"
	app2 "github.com/bighuangbee/basic-service/internal/module/operationLog/app"
	repo2 "github.com/bighuangbee/basic-service/internal/module/operationLog/repo"
	service2 "github.com/bighuangbee/basic-service/internal/module/operationLog/service"
	"github.com/bighuangbee/basic-service/internal/pkg/middleware"
	"github.com/bighuangbee/basic-service/internal/protocol"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "net/http/pprof"
)

// Injectors from wire.go:

func autoWireApp(bootstrap *conf.Bootstrap, logger log.Logger, helper *log.Helper, opLog *middleware.OpLog) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(bootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	iAccountRepo := repo.NewAccountRepo(dataData, helper, bootstrap)
	accountService := service.NewAccountService(iAccountRepo, logger, bootstrap)
	accountServer := app.NewAccountApp(accountService, helper)
	iOperationLogRepo := repo2.NewOperationLogRepo(dataData)
	operationLogService := service2.NewOperationLogService(iOperationLogRepo, logger)
	operationLogServer := app2.NewOperationLogApp(operationLogService, logger)
	pbServer := &protocol.PbServer{
		Account: accountServer,
		OpLog:   operationLogServer,
	}
	server := protocol.NewHTTPServer(bootstrap, logger, pbServer, dataData, opLog)
	grpcServer := protocol.NewGRPCServer(bootstrap, logger, pbServer)
	kratosApp := newApp(bootstrap, logger, server, grpcServer, opLog)
	return kratosApp, func() {
		cleanup()
	}, nil
}
