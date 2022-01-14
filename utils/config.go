package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	Cfg *viper.Viper
}

func LoadConfig() (config Config, err error) {
	// Config
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{Cfg: viper.GetViper()}, nil
}
func (c Config) Sub(path string) Config {
	return Config{Cfg: c.Cfg.Sub(path)}
}

func (c Config) Unmarshal(o interface{}) {
	c.Cfg.Unmarshal(o)
}
