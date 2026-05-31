package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DBURL        string `mapstructure:"DATAB_URL"`
	RedisURL     string `mapstructure:"REDIS_URL"`
	OpenAIKey    string `mapstructure:"OPENAI_API_KEY"`
	GeminiKey    string `mapstructure:"GEMINI_API_KEY"`
	AnthropicKey string `mapstructure:"ANTHROPIC_API_KEY"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: .env file not found, using enviroment variables.")
	}

	err = viper.Unmarshal(&config)
	return
}
