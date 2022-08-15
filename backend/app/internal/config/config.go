package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE" env-default:"port" env-description:"'port' or 'sock'. if 'sock' then env 'SOCKET_FILE' is required"`
		BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"PORT" env-default:"5005"`
		SocketFile string `env:"SOCKET_FILE" env-default:"app.sock"`
		ServerIP   string `env:"SERVER_IP" env-default:"http://videot4pe.dev"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL" env-default:"trace"`
		JwtSecret string `env:"JWT_SECRET" env-default:"secret"`
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-default:"admin"`
			Password string `env:"ADMIN_PWD" env-default:"admin"`
		}
	}
	PostgreSQL struct {
		Username string `env:"PGUSER" env-default:"postgres"`
		Host     string `env:"PGHOST" env-default:"postgres"`
		Password string `env:"PGPASSWORD" env-default:"123"`
		Database string `env:"PGDATABASE" env-default:"smer"`
		Port     string `env:"PGPORT" env-default:"5432"`
		Timeout  uint16 `env:"PGTIMEOUT" env-default:"5000"`
	}
	Mailer struct {
		Host     string `env:"MAILER_HOST" env-default:"smtp.gmail.com"`
		Port     int    `env:"MAILER_PORT" env-default:"587"`
		Username string `env:"MAILER_USERNAME" env-default:"email@email.email"`
		Password string `env:"MAILER_PASSWORD" env-default:"password"`
	}
	OAuth struct {
		Google struct {
			Key         string `env:"GOOGLE_OAUTH_KEY"`
			Secret      string `env:"GOOGLE_OAUTH_SECRET"`
			CallbackUrl string `env:"GOOGLE_OAUTH_CALLBACK_URL" env-default:"http://localhost:5005/api/oauth/google/callback"`
		}
		VK struct {
			Key         string `env:"VK_OAUTH_KEY"`
			Secret      string `env:"VK_OAUTH_SECRET"`
			CallbackUrl string `env:"VK_OAUTH_CALLBACK_URL" env-default:"http://localhost:5005/api/oauth/vk/callback"`
		}
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		// TODO fix path
		if err := cleanenv.ReadConfig("../../.env", instance); err != nil {
			var helpText string
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
