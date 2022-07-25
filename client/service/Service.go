package service

import (
	"fmt"

	"github.com/ability-sh/abi-db/client"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientService struct {
	config interface{}    `json:"-"`
	name   string         `json:"-"`
	Addr   string         `json:"addr"`
	client *client.Client `json:"-"`
}

func newClientService(name string, config interface{}) *ClientService {
	return &ClientService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *ClientService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *ClientService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *ClientService) OnInit(ctx micro.Context) error {

	dynamic.SetValue(s, s.config)

	conn, err := grpc.Dial(s.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}

	s.client = client.NewClient(conn)

	return nil
}

/**
* 校验服务是否可用
**/
func (s *ClientService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *ClientService) Recycle() {

}

func GetClient(ctx micro.Context, name string) (*client.Client, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(*ClientService)
	if ok {
		return ss.client, nil
	}
	return nil, fmt.Errorf("service %s not instanceof *ClientService", name)
}
