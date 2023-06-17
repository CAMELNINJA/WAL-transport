package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/CAMELNINGA/cdc-postgres.git/config"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/kafka"
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
	logger.Info("Starting copy deamon,  version: %s\n", Version)
	if baseConfig.IsKatka {
		logger.Info("Starting kafka producer")
		if err := kafkaRun(ctx, logger, baseConfig.Kafka); err != nil {
			logger.Fatal(err)
		}
	}
}

func kafkaRun(ctx context.Context, logger *logrus.Entry, cfg config.Kafka) error {
	logger.Info("Starting kafka producer")
	var b kafka.Bits
	b = kafka.Set(b, kafka.Consumer)
	kafka := kafka.NewKafka(
		kafka.WithBrokers(cfg.Brokers),
		kafka.WithTopic(cfg.Topic),
		kafka.WithFlags(b),
		kafka.WithGroupID(cfg.GroupID),
	)

	return nil
}
