package protocol

import (
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/basic-service/internal/data"
	kitKratos "github.com/bighuangbee/gokit/kratos"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"time"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(bc *conf.Bootstrap, logger log.Logger, server *PbServer, data *data.Data) *http.Server {
	// 不需要验证token的地址
	//checkTokenWhiteList := []string{
	//	"/api.mozi.device.v1.Device/SyncWvp",
	//	"/api.mozi.device.v1.Device/DeviceRecordHook",
	//}
	c := bc.Server
	srv := http.NewServer(
		http.Address(c.HTTP.Addr),
		http.Timeout(time.Duration(c.HTTP.Timeout)*time.Second),
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate(),
			//opLogCli.SaveOpLog(),
			//hiKratos.HTTPReturnTraceID(),
		),

		http.Logger(logger),
		kitKratos.SuccessEncoder(),
		kitKratos.ErrorEncoder(),
	)
	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	server.RegisterHTTP(srv)

	return srv
}
