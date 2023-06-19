package config

type Cli struct {
	Kafka   Kafka             `json:"kafka"`
	Deamons map[string]Config `json:"deamons"`
}
