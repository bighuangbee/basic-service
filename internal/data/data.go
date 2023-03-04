package data

import (
	"context"
	"github.com/bighuangbee/basic-service/internal/conf"
	"github.com/bighuangbee/gokit/storage/kitRedis"
	"github.com/bighuangbee/gokit/tools/id"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var (
	defaultRedisPrefix = ""
)

type Data struct {
	dbInfo      *conf.Database
	db          *gorm.DB
	rdb         kitRedis.Client
	id 			id.Generater
}

func NewData(bc *conf.Bootstrap, logger log.Logger) (*Data, func(), error) {
	c := bc.Data

	defaultRedisPrefix = c.Redis.AutoPrefix
	var db *gorm.DB
	//db, err := kitGorm.New(&kitGorm.Options{
	//	Address:  c.Database.Address,
	//	UserName: c.Database.UserName,
	//	Password: c.Database.Password,
	//	DBName:   c.Database.DBName,
	//	// Tracer:   otel.GetTracerProvider(),
	//	Logger:  kitGorm.Logger{L: log.NewHelper(logger)},
	//	Charset: "utf8mb4",
	//})
	//if err != nil {
	//	return nil, nil, err
	//}

	sfId, err := id.New(bc.Server.NodeId)
	if err != nil {
		panic("snowflakeId fail" + err.Error())
	}

	//rClient, err := kitRedis.New(&kitRedis.Options{
	//	Addr:     c.Redis.Address,
	//	Password: c.Redis.Password,
	//	DB:       int(c.Redis.DB),
	//})
	//if err != nil {
	//	return nil, nil, err
	//}

	logger.Log(log.LevelDebug, "db connect:", c.Database.Address, ",driver:", c.Database.Driver)

	d := &Data{
		dbInfo:      c.Database,
		db:          db,
		//rdb:         rClient,
		id: sfId,
	}

	return d, func() {
		d.rdb.Close()
	}, nil
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *Data) Redis(prefix string) *kitRedis.MyRedis {
	return &kitRedis.MyRedis{Prefix: defaultRedisPrefix + prefix, Rdb: d.rdb}
}

func (d *Data) Id() id.Generater {
	return d.id
}
