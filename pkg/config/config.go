package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const envPrefix = "boo"

type Config struct {
	Port          int        `json:"port" yaml:"port"`
	Dsn           string     `json:"dsn" yaml:"dsn"`
	TgSecret      string     `json:"tgSecret" yaml:"tgSecret"`
	InitBotSecret string     `json:"initBotSecret" yaml:"initBotSecret"`
	Providers     []Provider `json:"providers" yaml:"providers"`
}

type Provider struct {
	Name     string `json:"name" yaml:"name"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Secret   string `json:"secret" yaml:"secret"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(v *viper.Viper) error {
	configFile := v.GetString("config")
	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file, %w", err)
		}
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	if err := v.Unmarshal(&c); err != nil {
		return fmt.Errorf("error unmarshalling config, %w", err)
	}
	return nil
}
