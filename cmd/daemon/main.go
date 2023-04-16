package main

// import (
// 	"context"
// 	"encoding/binary"
// 	"os"
// 	"os/signal"

// 	"github.com/CAMELNINGA/cdc-postgres.git/config"
// 	"github.com/CAMELNINGA/cdc-postgres.git/internal/parser"
// 	"github.com/CAMELNINGA/cdc-postgres.git/internal/repository"
// 	"github.com/CAMELNINGA/cdc-postgres.git/internal/usecase"
// 	"github.com/CAMELNINGA/cdc-postgres.git/pkg/postgres"
// 	"github.com/sirupsen/logrus"
// )

// // / logger log levels.
// const (
// 	warningLoggerLevel = "warning"
// 	errorLoggerLevel   = "error"
// 	fatalLoggerLevel   = "fatal"
// 	infoLoggerLevel    = "info"
// )

// // initLogger init logrus preferences.
// func initLogger(cfg config.LoggerCfg, version string) *logrus.Entry {
// 	logger := logrus.New()

// 	logger.SetReportCaller(cfg.Caller)

// 	if cfg.Format == "json" {
// 		logger.SetFormatter(&logrus.JSONFormatter{})
// 	}

// 	var level logrus.Level

// 	switch cfg.Level {
// 	case warningLoggerLevel:
// 		level = logrus.WarnLevel
// 	case errorLoggerLevel:
// 		level = logrus.ErrorLevel
// 	case fatalLoggerLevel:
// 		level = logrus.FatalLevel
// 	case infoLoggerLevel:
// 		level = logrus.InfoLevel
// 	default:
// 		level = logrus.DebugLevel
// 	}

// 	logger.SetLevel(level)

// 	return logger.WithField("version", version)
// }

// func main() {
// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer cancel()

// 	config := &config.Config{
// 		Database: postgres.DatabaseCfg{
// 			Host:     "localhost",
// 			Port:     5433,
// 			Name:     "postgres",
// 			User:     "postgres",
// 			Password: "pass",
// 		},
// 		LoggerCfg: config.LoggerCfg{
// 			Caller: false,
// 			Level:  "info",
// 			Format: "json",
// 		},
// 		Listener: config.Listener{
// 			RefreshConnection: 30,
// 		},
// 	}
// 	logger := initLogger(config.LoggerCfg, "1.0.0")
// 	pgConf := config.Database
// 	conn, rConn, err := postgres.InitPgxConnections(pgConf)
// 	if err != nil {
// 		logger.Error(err)
// 		return
// 	}

// 	service := usecase.NewWalListener(
// 		logger,
// 		"test",
// 		repository.NewRepository(conn),
// 		rConn,
// 		parser.NewBinaryParser(binary.BigEndian),
// 	)
// 	if err := service.Process(ctx); err != nil {
// 		logger.Error("service process: %w", err)
// 		return
// 	}
// }
