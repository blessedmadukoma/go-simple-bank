package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
type Config struct {
	DB_DRIVER      string `mapstructure:"DB_DRIVER"`
	DB_SOURCE      string `mapstructure:"DB_SOURCE"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}