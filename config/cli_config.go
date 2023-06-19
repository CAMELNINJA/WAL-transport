package config

import (
	"encoding/json"
	"os"
)

type Cli struct {
	Kafka   Kafka             `json:"kafka"`
	Deamons map[string]Config `json:"deamons"`
}

func (c *Cli) Parse(f *os.File) error {
	return json.NewDecoder(f).Decode(c)

}
