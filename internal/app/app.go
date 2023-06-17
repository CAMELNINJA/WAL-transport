package app

import (
	"context"
	"encoding/binary"

	"github.com/CAMELNINGA/cdc-postgres.git/config"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/kafka"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/parser"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/repository"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/usecase"
	"github.com/CAMELNINGA/cdc-postgres.git/pkg/postgres"
	"github.com/sirupsen/logrus"
)

func RunCopyDeamon(ctx context.Context, logger *logrus.Entry, conf *config.Config) error {
	conn, rConn, err := postgres.InitPgxConnections(conf.Database)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer conn.Close()
	defer rConn.Close()
	var b kafka.Bits
	b = kafka.Set(b, kafka.Producer)
	kafka := kafka.NewKafka(
		kafka.WithBrokers(conf.Kafka.Brokers),
		kafka.WithTopic(conf.Kafka.Topic),
		kafka.WithFlags(b),
	)
	service := usecase.NewWalListener(
		logger,
		conf.Listener.SlotName,
		repository.NewRepository(conn),
		rConn,
		parser.NewBinaryParser(binary.BigEndian),
		kafka,
	)
	if err := service.Process(ctx); err != nil {
		logger.Error("service process: %w", err)
		return err
	}
	return nil
}

func RunSaveDeamon(logger *logrus.Entry, conf *config.Config) error {

	return nil
}
