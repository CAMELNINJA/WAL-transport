package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
	erorwallistner "github.com/CAMELNINGA/cdc-postgres.git/pkg/error_walListner"
)

type usecase interface {
	SaveData(ctx context.Context, messages models.Message) error
}

type config interface {
	NewConfig() error
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

func (k kafk) ListenConfig(ctx context.Context, ctxconfig config) error {
	if k.consumer == nil {
		return erorwallistner.ErrNotConnectedKafaConsumer
	}
	for {
		msg, err := k.consumer.ReadMessage(k.ctx)
		if err != nil {
			return err
		}

		var messages config.Config
		if err := json.Unmarshal(msg.Value, &messages); err != nil {
			fmt.Println(err)
			continue
		}

		if err := config.NewConfig(); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (k kafk) Close() error {
	if k.consumer == nil {
		return erorwallistner.ErrNotConnectedKafaConsumer
	}
	return k.consumer.Close()
}
