package main

import (
	"flag"
	"github.com/bighuangbee/basic-service/internal/conf"
	kitLog "github.com/bighuangbee/gokit/log"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap/zapcore"
	grpcDial "google.golang.org/grpc"
	_ "net/http/pprof"
	"time"
)

var flagConf string

func newApp(bc *conf.Bootstrap, logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	var registrar registry.Registrar
	if bc.Discovery.OnOff {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   bc.MicroService.Etcd.Addr,
			DialTimeout: time.Second, DialOptions: []grpcDial.DialOption{grpcDial.WithBlock()},
		})
		if err != nil {
			panic(err)
		}

		registrar = etcd.New(client)
	}

	return kratos.New(
		kratos.Name(bc.Name),
		kratos.Version(bc.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(hs, gs),
		kratos.Registrar(registrar),
	)
}

func main() {
	flag.StringVar(&flagConf, "conf", "../../config/config.dev.yaml", "config path, eg: -conf config.yaml")
	flag.Parse()

	var bc *conf.Bootstrap
	c := config.New(config.WithSource(file.NewSource(flagConf)))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	zapLog := kitLog.NewZapLogger(&kitLog.Options{
		Level:       zapcore.DebugLevel,
		Skip:        2,
		Writer:      kitLog.NewFileWriter(&kitLog.FileOption{
			Filename: "/opt/logs/account-svc/%Y-%m-%d.log",
			MaxSize:  20,
		}),
	})
	logger := log.With(zapLog, "tid", tracing.TraceID())

	logHelper := log.NewHelper(logger)

	app, cleanup, err := autoWireApp(bc, logger, logHelper)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
