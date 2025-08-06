package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool   `yaml:"is_debug" env-required:"true"`
	Listen  Listen  `yaml:"listen"`
	Storage Storage `yaml:"storage"`
	Log     Log     `yaml:"log"`
}

type Listen struct {
	Host         string        `yaml:"host" env-required:"true"`
	BindIP       string        `yaml:"bind_ip" env-required:"true"`
	Port         string        `yaml:"port" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-required:"true"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-required:"true"`
	IDLETimeout  time.Duration `yaml:"idle_timeout" env-required:"true"`
	DBCtxTimeout time.Duration `yaml:"db_ctx_timeout" env-required:"true"`
}

type Storage struct {
	Psql       Psql  `yaml:"psql"`
	MemoryRepo *bool `yaml:"memory_repo" env-required:"true"`
}

type Psql struct {
	Host                         string `yaml:"host" env-required:"true"`
	Port                         string `yaml:"port" env-required:"true"`
	Database                     string `yaml:"database" env-required:"true"`
	Username                     string `yaml:"username" env-required:"true"`
	Password                     string `yaml:"password" env-required:"true"`
	MaxConnectionPoolSize        int32  `yaml:"max_connection_pool_size" env-required:"true"`
	MaxConnectionLifetimeMinutes int    `yaml:"max_connection_lifetime_minutes" env-required:"true"`
	ConnectionAcquisitionTimeout int    `yaml:"connection_acquisition_timeout" env-required:"true"`
}

type Redis struct {
	Host                         string `yaml:"host" env-required:"true"`
	Port                         string `yaml:"port" env-required:"true"`
	Database                     string `yaml:"database" env-required:"true"`
	Username                     string `yaml:"username" env-required:"true"`
	Password                     string `yaml:"password" env-required:"true"`
	MaxConnectionPoolSize        int32  `yaml:"max_connection_pool_size" env-required:"true"`
	MaxConnectionLifetimeMinutes int    `yaml:"max_connection_lifetime_minutes" env-required:"true"`
	ConnectionAcquisitionTimeout int    `yaml:"connection_acquisition_timeout" env-required:"true"`
}

type Log struct {
	Path     string `yaml:"path" env-required:"true"`
	Filename string `yaml:"filename" env-required:"true"`
}

var instance *Config
var once sync.Once

func Init() *Config {
	once.Do(func() {

		pwd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		pathConfig := pwd + "/config.yml"

		log.Println("read application configuration: pwd: ", pwd)

		instance = &Config{}
		if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
