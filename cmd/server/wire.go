// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/basic-service/internal/data"
	"github.com/bighuangbee/basic-service/internal/module/account"
	"github.com/bighuangbee/basic-service/internal/protocol"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func autoWireApp(*conf.Bootstrap, log.Logger, *log.Helper) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, protocol.ProviderSet, account.ProviderSet, newApp))
}
