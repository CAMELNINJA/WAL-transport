package config

import (
	"fmt"
	"time"

	"github.com/CAMELNINGA/cdc-postgres/pkg/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	Database postgres.DatabaseCfg `mapstructure:"db"`
	Listener Listener             `mapstructure:"listener"`
}

type Listener struct {
	RefreshConnection time.Duration
}

func NewConfig() *Config {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}
