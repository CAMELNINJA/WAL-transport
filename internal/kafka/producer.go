package kafka

import (
	"encoding/json"

	"github.com/CAMELNINGA/cdc-postgres.git/config"
	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
	error_walListner "github.com/CAMELNINGA/cdc-postgres.git/pkg/error_walListner"
	kafka "github.com/segmentio/kafka-go"
)

func (k *kafk) Publish(key string, messages models.Message) error {
	if k.producer == nil {
		return error_walListner.ErrNotConnectedKafaProducer
	}

	value, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return k.producer.WriteMessages(k.ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

func (k *kafk) PublishConfig(key string, messages config.Config) error {
	if k.producer == nil {
		return error_walListner.ErrNotConnectedKafaProducer
	}
	value, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return k.producer.WriteMessages(k.ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}
