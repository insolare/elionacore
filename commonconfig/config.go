package commonconfig

import (
	"fmt"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type BrokerList []string
type Config struct {
	DSN   string     `env:"CONNECTION_STRING" env-required:"true"`
	Seeds BrokerList `env:"BROKERS" env-required:"true"`
}

func GetCommonConfig() (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)

	return cfg, err
}

func (b *BrokerList) SetValue(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("no brokers provided")
	}

	s = strings.ReplaceAll(s, "|", ",")
	brokers := strings.Split(s, ",")
	*b = brokers

	return nil
}
