package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

var AppConfig = loadConfig()

type Config struct {
	TgApiUrl     string `yaml:"TG_URL" env:"TG_URL"`
	TgApiToken   string `yaml:"TG_TOKEN" env:"TG_TOKEN"`
	TemplatePath string `yaml:"TEMPLATE_PATH" env:"TEMPLATE_PATH" env-default:"./default-template.txt"`
	IqairUrl     string `yaml:"IQAIR_URL" env:"IQAIR_URL" env-default:"http://api.airvisual.com/v2"`
	IqairToken   string `yaml:"IQAIR_TOKEN" env:"IQAIR_TOKEN"`
	Country      string `yaml:"COUNTRY" env:"COUNTRY" env-default:"Russia"`
	State        string `yaml:"STATE" env:"STATE" env-default:"Moscow"`
	City         string `yaml:"CITY" env:"CITY" env-default:"Moscow"`
	HttpPort     int    `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8080"`
	LogLevel     string `yaml:"LOG_LEVEL" env:"LOG_LEVEL" env-default:"info"`
}

func loadConfig() *Config {
	var cfg Config
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("Unable to ReadEnv: %v", err)
		}
	} else {
		err = cleanenv.ReadConfig("config.yaml", &cfg)
		if err != nil {
			log.Fatalf("Unable to ReadConfig: %v", err)
		}
	}

	cfg.initLogger()

	return &cfg
}

func (c *Config) initLogger() {
	const defaultLogLevel = "info"

	logLevelValue, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Warnf(
			"Invalid 'LOG_LEVEL': [%v], will use default LOG_LEVEL: [%v]. "+
				"Allowed values: trace, debug, info, warn, warning, error, fatal, panic", c.LogLevel, defaultLogLevel,
		)

		logLevelValue = log.InfoLevel
	}
	log.SetLevel(logLevelValue)
}

