package app

import (
	"context"
	"encoding/binary"

	"github.com/CAMELNINGA/WAL-transport.git/config"
	"github.com/CAMELNINGA/WAL-transport.git/internal/kafka"
	"github.com/CAMELNINGA/WAL-transport.git/internal/parser"
	querybuilder "github.com/CAMELNINGA/WAL-transport.git/internal/query_builder"
	"github.com/CAMELNINGA/WAL-transport.git/internal/repository"
	"github.com/CAMELNINGA/WAL-transport.git/internal/sanitize"
	"github.com/CAMELNINGA/WAL-transport.git/internal/usecase"
	"github.com/CAMELNINGA/WAL-transport.git/pkg/postgres"
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
		kafka.WithContext(ctx),
	)
	sant, err := initSanitase(logger, conf)
	if err != nil {
		logger.Error(err)
		return err
	}

	service := usecase.NewWalListener(
		logger,
		conf.Listener.SlotName,
		repository.NewRepository(conn),
		rConn,
		parser.NewBinaryParser(binary.BigEndian),
		kafka,
		sant,
	)
	if err := service.Process(ctx); err != nil {
		logger.Error("service process: %w", err)
		return err
	}
	return nil
}

func RunSaveDeamon(ctx context.Context, logger *logrus.Entry, conf *config.Config) error {
	masterConn, err := postgres.InitMasterConnection(ctx, conf.Database)
	if err != nil {
		logger.Error(err)
		return err
	}
	sanit, err := initSanitase(logger, conf)
	if err != nil {
		logger.Error(err)
		return err
	}
	collector := usecase.NewCollector(querybuilder.NewQueryBuilder(logger), masterConn, sanit)
	var b kafka.Bits
	b = kafka.Set(b, kafka.Consumer)
	kafka := kafka.NewKafka(
		kafka.WithBrokers(conf.Kafka.Brokers),
		kafka.WithTopic(conf.Kafka.Topic),
		kafka.WithFlags(b),

		kafka.WithGroupID(conf.Kafka.GroupID),
	)

	return kafka.Listen(ctx, collector)
}

func initSanitase(logger *logrus.Entry, conf *config.Config) (sanitize.Handler, error) {
	sant := sanitize.NewSanitizeHandler()
	handlers := []sanitize.Handler{}
	handlers = append(handlers, sant)
	for _, s := range conf.Sanitize {
		switch s.Type {
		case config.FilterType:
			logger.Info("Init filter")
			opts := []sanitize.FilterOpts{}

			if s.Columns != nil {
				opts = append(opts, sanitize.WithFilterColumns(s.Columns))
			}
			if s.Table != "" {
				opts = append(opts, sanitize.WithFilterTable(s.Table))
			}
			if s.Schema != nil {
				opts = append(opts, sanitize.WithFilterSchema(s.Schema))
			}
			filter := sanitize.NewFilterHandler(opts...)
			handlers = append(handlers, filter)
		case config.ReplaseType:
			logger.Info("Init replase")
			opts := []sanitize.ReplaceOpts{}

			if s.Columns != nil {
				opts = append(opts, sanitize.WithReplaceColumns(s.Columns))
			}
			if s.Table != "" && s.OldTable != "" {
				opts = append(opts, sanitize.WithReplaceTable(s.Table, s.OldTable))
			}
			if s.Schema != nil {
				opts = append(opts, sanitize.WithReplaceSchema(s.Schema))
			}
			replase := sanitize.NewReplaceHandler(opts...)
			handlers = append(handlers, replase)

		default:
			logger.Info("Unknown type")
		}

	}
	for i, h := range handlers {
		if i == len(handlers)-1 {
			break
		}
		h.SetNext(handlers[i+1])
	}

	return sant, nil
}

func KafkaRun(ctx context.Context, logger *logrus.Entry, cfg config.Kafka, cfgChan chan config.Config) error {
	logger.Info("Starting kafka producer")
	var b kafka.Bits
	b = kafka.Set(b, kafka.Consumer)
	kafka := kafka.NewKafka(
		kafka.WithBrokers(cfg.Brokers),
		kafka.WithTopic(cfg.Topic),
		kafka.WithFlags(b),
		kafka.WithGroupID(cfg.GroupID),
	)

	return kafka.ListenConfig(ctx, cfgChan)

}
