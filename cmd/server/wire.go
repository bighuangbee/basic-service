// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/basic-service/internal/data"
	"github.com/bighuangbee/basic-service/internal/module/account"
	"github.com/bighuangbee/basic-service/internal/module/operationLog"
	"github.com/bighuangbee/basic-service/internal/pkg/middleware"
	"github.com/bighuangbee/basic-service/internal/protocol"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func autoWireApp(*conf.Bootstrap, log.Logger, *log.Helper, *middleware.OpLog) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, protocol.ProviderSet, newApp,
		account.ProviderSet,
		operationLog.ProviderSet,
	))
}
