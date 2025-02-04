package config

type Config struct {
	Listen    int        `json:"listen" yaml:"listen"`
	Dsn       string     `json:"dsn" yaml:"dsn"`
	TgSecret  string     `json:"tgSecret" yaml:"tgSecret"`
	Providers []Provider `json:"providers" yaml:"providers"`
}

type Provider struct {
	Name     string `json:"name" yaml:"name"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Secret   string `json:"secret" yaml:"secret"`
}

func NewConfig() *Config {
	return &Config{}
}
