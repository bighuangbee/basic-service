package conf

type MicroService struct {
	OpLog    *ServiceInfo
	Etcd    Etcd
}

type ServiceInfo struct {
	Name    string
	Timeout int32
	Grpc    string
	IsLocal bool
}

type Etcd struct {
	Addr []string
}
