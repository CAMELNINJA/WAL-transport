package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CAMELNINGA/WAL-transport.git/config"
	"github.com/CAMELNINGA/WAL-transport.git/internal/models"
	erorwallistner "github.com/CAMELNINGA/WAL-transport.git/pkg/error_walListner"
)

type usecase interface {
	SaveData(ctx context.Context, messages models.Message) error
}

func (k kafk) Listen(ctx context.Context, usecase usecase) error {
	if k.consumer == nil {
		return erorwallistner.ErrNotConnectedKafaConsumer
	}
	for {
		msg, err := k.consumer.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var messages models.Message
		if err := json.Unmarshal(msg.Value, &messages); err != nil {
			fmt.Println(err)
			continue
		}

		if err := usecase.SaveData(ctx, messages); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (k kafk) ListenConfig(ctx context.Context, cfg chan<- config.Config) error {
	if k.consumer == nil {
		return erorwallistner.ErrNotConnectedKafaConsumer
	}

	for {
		msg, err := k.consumer.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if string(msg.Key) != k.GroupID {
			continue
		}
		var messages config.Config
		if err := json.Unmarshal(msg.Value, &messages); err != nil {
			fmt.Println(err)
			continue
		}

		cfg <- messages
	}
}

func (k kafk) Close() error {
	if k.consumer == nil {
		return erorwallistner.ErrNotConnectedKafaConsumer
	}
	return k.consumer.Close()
}
