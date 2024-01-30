package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment                  string   `mapstructure:"ENVIRONMENT"`
	DBDriver                     string   `mapstructure:"DB_DRIVER"`
	DBSource                     string   `mapstructure:"DB_SOURCE"`
	RedisConnURL                 string   `mapstructure:"REDIS_CONN_URL"`
	DBMigrationURL               string   `mapstructure:"DB_MIGRATION_URL"`
	HTTPServerAddress            string   `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress            string   `mapstructure:"GRPC_SERVER_ADDRESS"`
	AuthServiceGRPCServerAddress string   `mapstructure:"AUTH_SERVICE_GRPC_SERVER_ADDRESS"`
	ServiceAuthPublicKeys        []string `mapstructure:"SERVICE_AUTH_PUBLIC_KEYS"`
	ServiceAuthPrivateKeys       []string `mapstructure:"SERVICE_AUTH_PRIVATE_KEYS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetTypeByDefaultValue(true)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
