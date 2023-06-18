package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/CAMELNINGA/cdc-postgres.git/config"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/app"
	"github.com/sirupsen/logrus"
)

var Version string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	baseConfig, err := config.NewBaseConfig()
	if err != nil {
		panic(err)
	}
	logger := config.InitLogger(baseConfig.LoggerCfg, Version)
	logger.Info(fmt.Printf("Starting copy deamon,  version: %s\n", Version))
	cfgChan := make(chan config.Config)
	defer close(cfgChan)
	if baseConfig.IsKatka {
		logger.Info("Starting kafka producer")
		if err := app.KafkaRun(ctx, logger, baseConfig.Kafka, cfgChan); err != nil {
			logger.Fatal(err)
		}
	}
}

func run(ctx context.Context, logger *logrus.Entry, <-cfg chan config.Config) error {
	logger.Info("Starting copy deamon")
	
	for {
		select {
		case <-ctx.Done():
			return nil
		case cfg := <-cfg:
			copyCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			if err := app.RunCopyDeamon(copyCtx,logger,cfg); err != nil {
				return err
			}

		}
	}
	
}
