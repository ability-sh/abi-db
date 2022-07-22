package srv

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("abi-db", func(name string, config interface{}) (micro.Service, error) {
		return newDBService(name, config), nil
	})
}
