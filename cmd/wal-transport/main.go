package main

import (
	"context"
	"fmt"

	"github.com/CAMELNINGA/cdc-postgres.git/config"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/kafka"
	"github.com/CAMELNINGA/cdc-postgres.git/pkg/postgres"
)

func main() {
	ctx := context.Background()
	conf := config.Config{
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
	}
	fmt.Println(conf.Database)

	saveconf := config.Config{
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
	}
	var flag kafka.Bits
	flag = kafka.Set(flag, kafka.Producer)
	k := kafka.NewKafka(
		kafka.WithBrokers(conf.Kafka.Brokers),
		kafka.WithTopic("Config"),
		kafka.WithGroupID("cli"),
		kafka.WithContext(ctx),
		kafka.WithFlags(flag),
	)
	k.PublishConfig("copy_deamon", conf)
	k.PublishConfig("save_deamon", saveconf)

}
