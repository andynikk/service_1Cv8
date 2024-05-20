package environment

import (
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v7"
)

// configENV структура для хранения параметров сеанса
type configENV struct {
	DatabaseDsn string `env:"DATABASE_URI"`
}

// DBConfig структура хранения свойств базы данных
type DBConfig struct {
	DatabaseDsn string
	Key         string
	Path        string
}

// Config структура хранения свойств конфигурации сервера
type Config struct {
	DBConfig
}

// NewConfig создание и заполнение структуры свойств сервера
func NewConfig() (*Config, error) {

	keyDatabaseDsn := flag.String("d", "", "строка соединения с базой")
	flag.Parse()

	var cfgENV configENV
	err := env.Parse(&cfgENV)
	if err != nil {
		log.Fatal(err)
	}

	databaseDsn := cfgENV.DatabaseDsn
	if _, ok := os.LookupEnv("DATABASE_URI"); !ok {
		databaseDsn = *keyDatabaseDsn
	}

	sc := Config{
		DBConfig: DBConfig{
			DatabaseDsn: databaseDsn,
		},
	}

	return &sc, err
}

func ServerUrlApi() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	port := constants.Port

	data := repository.DataJSON{}
	_ = repository.GetPudgeSetting(&data.Settings)
	if data.Settings.HTTPPort != "" {
		port = data.Settings.HTTPPort
	}

	return fmt.Sprintf("http://%s:%s/api", hostname, port)
}
