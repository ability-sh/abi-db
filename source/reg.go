package source

import "fmt"

type Registry = func(driver string, config interface{}) (Source, error)

var sourceSet = map[string]Registry{}

func Reg(dirver string, registry Registry) {
	sourceSet[dirver] = registry
}

func NewSource(driver string, config interface{}) (Source, error) {
	r, ok := sourceSet[driver]
	if !ok {
		return nil, fmt.Errorf("not found %s registry", driver)
	}
	return r(driver, config)
}
