package kafka

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type Bits uint8

const (
	Producer Bits = 1 << iota
	Consumer
)

func Set(b, flag Bits) Bits    { return b | flag }
func Clear(b, flag Bits) Bits  { return b &^ flag }
func Toggle(b, flag Bits) Bits { return b ^ flag }
func Has(b, flag Bits) bool    { return b&flag != 0 }

type kafk struct {
	producer *kafka.Writer
	consumer *kafka.Reader
	ctx      context.Context
	Brokers  []string
	Topic    string
	GroupID  string
	flags    Bits
}

type KafaOption func(*kafk)

func WithBrokers(brokers []string) KafaOption {
	return func(k *kafk) {
		k.Brokers = brokers
	}
}

func WithTopic(topic string) KafaOption {
	return func(k *kafk) {
		k.Topic = topic
	}
}
func WithFlags(flags Bits) KafaOption {
	return func(k *kafk) {
		k.flags = flags
	}
}

func WithGroupID(groupID string) KafaOption {
	return func(k *kafk) {
		k.GroupID = groupID
	}
}

func WithContext(ctx context.Context) KafaOption {
	return func(k *kafk) {
		k.ctx = ctx
	}
}

func NewKafka(opts ...KafaOption) *kafk {
	k := &kafk{}
	for _, opt := range opts {
		opt(k)
	}
	if Has(k.flags, Producer) {
		k.producer = kafka.NewWriter(kafka.WriterConfig{
			Brokers: k.Brokers,
			Topic:   k.Topic,
		})
	}
	if Has(k.flags, Consumer) {
		k.consumer = kafka.NewReader(kafka.ReaderConfig{
			Brokers: k.Brokers,
			Topic:   k.Topic,
			GroupID: k.GroupID,
		})
	}

	return k
}
