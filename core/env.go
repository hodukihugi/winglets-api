package core

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerPort                 string        `mapstructure:"SERVER_PORT"`
	Environment                string        `mapstructure:"ENV"`
	LogOutput                  string        `mapstructure:"LOG_OUTPUT"`
	LogLevel                   string        `mapstructure:"LOG_LEVEL"`
	DBUsername                 string        `mapstructure:"DB_USER"`
	DBPassword                 string        `mapstructure:"DB_PASS"`
	DBHost                     string        `mapstructure:"DB_HOST"`
	DBPort                     string        `mapstructure:"DB_PORT"`
	DBName                     string        `mapstructure:"DB_NAME"`
	SmtpUser                   string        `mapstructure:"SMTP_USER"`
	SmtpPassword               string        `mapstructure:"SMTP_PASS"`
	SmtpHost                   string        `mapstructure:"SMTP_HOST"`
	JWTSecret                  string        `mapstructure:"JWT_SECRET"`
	AccessTokenExpiresIn       time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn      time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	EmailVerificationExpiresIn time.Duration `mapstructure:"EMAIL_VERIFICATION_EXPIRED_IN"`
}

// NewEnv creates a new environment
func NewEnv() *Env {

	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	return &env
}
