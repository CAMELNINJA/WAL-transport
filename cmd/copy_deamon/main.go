package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/CAMELNINGA/WAL-transport.git/config"
	"github.com/CAMELNINGA/WAL-transport.git/internal/app"
	"github.com/sirupsen/logrus"
)

var Version string = "1.0.0"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	baseConfig, err := config.NewBaseConfig()
	if err != nil {
		panic(err)
	}
	logger := config.InitLogger(baseConfig.LoggerCfg, Version)
	logger.Info(fmt.Printf("Starting copy deamon,  version: %s\n", Version))
	shutdown := make(chan error, 1)
	defer close(shutdown)
	cfgChan := make(chan config.Config)

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
			cfg, ok := <-cfgChan
			if !ok {
				logger.Info("stop config listener")
				return
			}
			fmt.Println("newCfg", newCfg)
			if newCfg {
				logger.Info("New config received")
				stop <- struct{}{}
			}
			fmt.Println(cfg)
			go func(stop <-chan struct{}) {
				newCfg = true
				shutdown <- run(stop, logger, &cfg)
			}(stop)

		}
	}(stop)

	select {
	case <-ctx.Done():
		stop <- struct{}{}
		logger.Info("Shutdown signal received")
		close(cfgChan)
		cancel()
	case err := <-shutdown:
		stop <- struct{}{}
		logger.Error(err)
		close(cfgChan)
		cancel()
	}
	time.Sleep(2 * time.Second)
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
