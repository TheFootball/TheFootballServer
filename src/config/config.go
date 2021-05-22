package config

import (
	"log"

	"github.com/spf13/viper"
)

func GetEnv(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Printf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion %s", value)
	}

	return value
}

//var dsnConfigs = []string{
//	fmt.Sprintf("host=%s", GetEnv("POSTGRES_HOST")),
//	fmt.Sprintf("user=%s", GetEnv("POSTGRES_USER")),
//	fmt.Sprintf("database=%s", GetEnv("POSTGRES_DB")),
//	fmt.Sprintf("password=%s", GetEnv("POSTGRES_PASSWORD")),
//	fmt.Sprintf("port=%s", GetEnv("POSTGRES_PORT")),
//	"sslmode=disable",
//	fmt.Sprintf("TimeZone=%s", GetEnv("POSTGRES_TZ")),
//}

//var DSN string = strings.Join(dsnConfigs, " ")
