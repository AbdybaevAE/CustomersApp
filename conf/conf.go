package conf

import "github.com/spf13/viper"

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DbUser        string `mapstructure:"POSTGRES_USER"`
	DbPassword    string `mapstructure:"POSTGRES_PASSWORD"`
	DbName        string `mapstructure:"POSTGRES_DB"`
	DbHost        string `mapstructure:"POSTGRES_HOST"`
}

func Load() *Config {
	conf := &Config{}
	viper.AddConfigPath("./resources/")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("cannot read config from app.env file " + err.Error())
	}
	if err := viper.Unmarshal(conf); err != nil {
		panic("error unmarsha config")
	}
	return conf
}
