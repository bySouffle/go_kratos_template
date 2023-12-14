package cronjob

import "google.golang.org/protobuf/types/known/durationpb"

// CronConfig	CronRedis配置
type CronConfig struct {
	Network      string               `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Addr         string               `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Db           int32                `protobuf:"varint,3,opt,name=db,proto3" json:"db,omitempty"`
	Password     string               `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	DialTimeout  *durationpb.Duration `protobuf:"bytes,5,opt,name=dial_timeout,json=dialTimeout,proto3" json:"dial_timeout,omitempty"`
	ReadTimeout  *durationpb.Duration `protobuf:"bytes,6,opt,name=read_timeout,json=readTimeout,proto3" json:"read_timeout,omitempty"`
	WriteTimeout *durationpb.Duration `protobuf:"bytes,7,opt,name=write_timeout,json=writeTimeout,proto3" json:"write_timeout,omitempty"`
	MinIdleConn  int32                `protobuf:"varint,8,opt,name=MinIdleConn,proto3" json:"MinIdleConn,omitempty"`
	PoolSize     int32                `protobuf:"varint,9,opt,name=PoolSize,proto3" json:"PoolSize,omitempty"`
	PoolTimeout  *durationpb.Duration `protobuf:"bytes,10,opt,name=PoolTimeout,proto3" json:"PoolTimeout,omitempty"`
	Concurrency  int32                `protobuf:"varint,11,opt,name=Concurrency,proto3" json:"Concurrency,omitempty"`
}
