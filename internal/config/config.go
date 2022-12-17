package config

import (
	"flag"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Proxy    Proxy     `mapstructure:"proxy" json:"proxy"`
	Backends []Backend `mapstructure:"backends" json:"backends"`
}

type Proxy struct {
	Port string `mapstructure:"port" json:"port"`
}

type Backend struct {
	IsDead bool
	URL    string `mapstructure:"url" json:"url"`
	Mu     sync.RWMutex
}

func LoadConfigs() *Config {
	config := &Config{}
	var configFilepath string
	flag.StringVar(&configFilepath, "config file path", "config.yml", "Path to config file")
	v := viper.New()
	v.SetConfigFile(configFilepath)
	v.ReadInConfig()
	unmarshalConfig(config, v)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config has changed")
		unmarshalConfig(config, v)
	})
	return config
}

func unmarshalConfig(config *Config, v *viper.Viper) {
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("[CONFIG] Error unmarshaling app config : %+v\n", err)
	}
}
