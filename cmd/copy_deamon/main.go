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

var Version string = "1.0.0"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	baseConfig, err := config.NewBaseConfig()
	if err != nil {
		panic(err)
	}
	logger := config.InitLogger(baseConfig.LoggerCfg, Version)
	logger.Info(fmt.Printf("Starting copy deamon,  version: %s\n", Version))
	shutdown := make(chan error, 1)
	defer close(shutdown)
	cfgChan := make(chan config.Config)
	defer close(cfgChan)
	stop := make(chan struct{})
	defer close(stop)
	var newCfg bool

	if baseConfig.IsKatka {
		logger.Info("Starting kafka producer")

		go func(shutdown chan<- error, ctx context.Context) {

			shutdown <- app.KafkaRun(ctx, logger, baseConfig.Kafka, cfgChan)
		}(shutdown, ctx)
	}

	go func(stop chan struct{}) {
		for {
			select {
			case cfg := <-cfgChan:
				if newCfg {
					logger.Info("New config received")
					stop <- struct{}{}
				}
				go func(stop <-chan struct{}) {
					newCfg = true
					shutdown <- run(stop, logger, &cfg)
				}(stop)
			}
		}
	}(stop)
	select {
	case <-ctx.Done():
		stop <- struct{}{}
		logger.Info("Shutdown signal received")
		cancel()
	case err := <-shutdown:
		stop <- struct{}{}
		logger.Error(err)
		cancel()
	}
	logger.Info("Shutdown complete")
}

func run(stop <-chan struct{}, logger *logrus.Entry, cfg *config.Config) error {
	logger.Info("Starting copy deamon")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := app.RunCopyDeamon(ctx, logger, cfg); err != nil {
		return err
	}
	<-stop
	ctx.Done()
	return nil
}
