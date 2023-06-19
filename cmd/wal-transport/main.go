package main

import (
	"context"

	"github.com/CAMELNINGA/cdc-postgres.git/config"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/kafka"
	"github.com/CAMELNINGA/cdc-postgres.git/pkg/postgres"
)

// CLi chicken config and send to kafka
func main() {
	ctx := context.Background()

	cli := config.Cli{
		Kafka: config.Kafka{
			Brokers: []string{"camelninja.ru:9092"},
			Topic:   "Config",
			GroupID: "cli",
		},
		Deamons: map[string]config.Config{
			"copy_deamon": {
				Database: postgres.DatabaseCfg{
					Host:     "slave-pg",
					Port:     5432,
					User:     "postgres",
					Name:     "postgres",
					Password: "pass",
				},
				Listener: config.Listener{
					RefreshConnection: 10000000 * 100000000000,
					SlotName:          "test",
				},
				Kafka: config.Kafka{
					Brokers: []string{"camelninja.ru:9092"},
					Topic:   "Data",
					GroupID: "copy_deamon",
				},
				Sanitize: []config.Sanitize{
					{
						Type:  "filter",
						Table: "*",
						Columns: map[string]string{
							"id": "id",
						},
					},
				},
			},
			"save_deamon": {
				Database: postgres.DatabaseCfg{
					Host:     "master-pg",
					Port:     5432,
					User:     "postgres",
					Name:     "postgres",
					Password: "pass",
				},
				Listener: config.Listener{
					RefreshConnection: 10000000 * 100000000000,
					SlotName:          "test",
				},
				Kafka: config.Kafka{
					Brokers: []string{"camelninja.ru:9092"},
					Topic:   "Data",
					GroupID: "save_deamon",
				},
				Sanitize: []config.Sanitize{
					{
						Type:  "filter",
						Table: "*",
						Columns: map[string]string{
							"id": "id",
						},
					},
				},
			},
		},
	}

	var flag kafka.Bits
	flag = kafka.Set(flag, kafka.Producer)
	k := kafka.NewKafka(
		kafka.WithBrokers(cli.Kafka.Brokers),
		kafka.WithTopic(cli.Kafka.Topic),
		kafka.WithGroupID(cli.Kafka.GroupID),
		kafka.WithContext(ctx),
		kafka.WithFlags(flag),
	)
	for name, v := range cli.Deamons {
		k.PublishConfig(name, v)
	}
}
