package configuration

import (
	"errors"
	"os"

	"github.com/ReanSn0w/kincong/internal/utils"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidConfiguration = errors.New("invalid configuration")
)

// Load - загружает конфигурацию из файла
func Load(filename string) (Configuration, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	conf := Configuration{}
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		return nil, ErrInvalidConfiguration
	}

	return conf, nil
}

type Configuration []Connection

type Connection struct {
	Title  string        `yaml:"title"`
	Values []utils.Value `yaml:"values,omitempty"`
}
