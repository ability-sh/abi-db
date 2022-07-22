package srv

import (
	"fmt"

	"github.com/ability-sh/abi-db/db"
	"github.com/ability-sh/abi-db/source"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
)

const (
	SERVICE_ABI_DB = "abi-db"
)

type DBService struct {
	config interface{} `json:"-"`
	name   string      `json:"-"`
	Driver string      `json:"driver"`
	db     db.Database `json:"-"`
}

func newDBService(name string, config interface{}) *DBService {
	return &DBService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *DBService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *DBService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *DBService) OnInit(ctx micro.Context) error {

	dynamic.SetValue(s, s.config)

	ss, err := source.NewSource(s.Driver, s.config)

	if err != nil {
		return err
	}

	s.db = db.Open(ss)

	return nil
}

/**
* 校验服务是否可用
**/
func (s *DBService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *DBService) Recycle() {

}

func GetDBService(ctx micro.Context, name string) (*DBService, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(*DBService)
	if ok {
		return ss, nil
	}
	return nil, fmt.Errorf("service %s not instanceof *DBService", name)
}
