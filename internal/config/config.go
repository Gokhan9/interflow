package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DBURL        string `mapstructure:"DATABASE_URL"`
	RedisURL     string `mapstructure:"REDIS_URL"`
	OpenAIKey    string `mapstructure:"OPENAI_API_KEY"`
	GeminiKey    string `mapstructure:"GEMINI_API_KEY"`
	AnthropicKey string `mapstructure:"ANTHROPIC_API_KEY"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigName(".env") // Config dosyasının adını belirtir. Burada .env olarak belirtiyoruz, ancak viper config dosyalarını genellikle JSON, YAML veya TOML formatlarında bekler. .env dosyaları genellikle KEY=VALUE formatında olduğu için, viper'ın bu formatı doğrudan desteklemediğini unutmayın. Ancak, viper.SetConfigType("env") ile bu durumu aşabiliriz.
	viper.SetConfigType("env")  // .env dosyalarının KEY=VALUE formatında olduğunu belirtir. Bu sayede viper, .env dosyasını doğru şekilde okuyabilir.
	viper.AddConfigPath(".")    // Çalışma dizininde .env dosyasını arar

	viper.AutomaticEnv() // Ortam değişkenlerini otomatik olarak okur.

	err = viper.ReadInConfig() // .env dosyasını okumaya çalışır. Eğer dosya bulunamazsa veya okunamazsa, bu fonksiyon bir hata döndürür.
	if err != nil {
		log.Printf("Warning: .env file not found, using enviroment variables.")
	}

	err = viper.Unmarshal(&config) // Viper tarafından okunan konfigürasyon değerlerini Config struct'ına dönüştürür. Eğer yapılandırma değerleri struct ile uyumlu değilse veya başka bir sorun varsa, bu fonksiyon bir hata döndürür.
	return
}
