package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/CAMELNINGA/WAL-transport.git/config"
	"github.com/CAMELNINGA/WAL-transport.git/internal/kafka"
	error_walListner "github.com/CAMELNINGA/WAL-transport.git/pkg/error_walListner"
)

func CheckConfig(filepath string) (string, error) {
	_, err := getConf(filepath)
	if err != nil {
		return "", err
	}

	return "The file is parsed successfully...", nil
}

func SendConfig(ctx context.Context, filepath string) (string, error) {
	cli, err := getConf(filepath)
	if err != nil {
		return "", err
	}

	if len(cli.Kafka.Brokers) == 0 {
		return "", error_walListner.ErrKafkaBrokersNotSet
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

		if err := k.PublishConfig(name, v); err != nil {
			return "", error_walListner.ErrSendConfigToKafka
		}
	}
	return "The config sended successfully...", nil
}

func getConf(filepath string) (*config.Cli, error) {
	cli := config.Cli{}
	f, err := os.Open(filepath)
	if err != nil {
		return nil, error_walListner.ErrConfigFileNotFound
	}
	defer f.Close()
	fmt.Println("The File is opened successfully...")

	if err := cli.Parse(f); err != nil {
		return nil, error_walListner.ErrConfigFileNotParsed
	}

	return &cli, nil
}
