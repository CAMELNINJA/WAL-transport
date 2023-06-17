package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type BaseConfig struct {
	LoggerCfg LoggerCfg `namespace:"logger" group:"logger" env-namespace:"LOGGER" env-group:"LOGGER"`
	IsKatka   bool      `long:"is-kafka" env:"IS_KAFKA" description:"Is kafka"`
	Kafka     Kafka     `namespace:"kafka" group:"kafka" env-namespace:"KAFKA" env-group:"KAFKA"`
}

func NewBaseConfig() (*BaseConfig, error) {
	var config BaseConfig
	p := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)

	_, err := p.ParseArgs(os.Args)
	if err != nil {
		return nil, err
	}

	return &config, nil

}
