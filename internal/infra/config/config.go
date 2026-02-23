package config

import (
	"fmt"
	"os"
	"time"
)

// Config is config app
type Config struct {
	BotAntiSpam botAntiSpam
	Storage     Storage
}

// NewConfig return config app instance
func NewConfig() (*Config, error) {
	wlTags := map[string]struct{}{
		"@prolann": {},
		"@Prolann": {},
	}

	conf := &Config{
		BotAntiSpam: botAntiSpam{
			Settings: settings{
				Token:                 os.Getenv("antispam_bot_token"),
				AdmToken:              os.Getenv("adm_antispam_bot_token"),
				OffsetMessageStart:    0,
				TimeOut:               60,
				Reties:                3,
				TimeOutBetweenRetries: 10 * time.Second,
			},
			WhiteListTags:   wlTags,
			WhiteListAuthor: []int64{445149872, 101316726},
		},
		Storage: Storage{
			hostDB:     os.Getenv("DB_HOST"),
			portDB:     os.Getenv("DB_PORT"),
			nameDB:     os.Getenv("DB_NAME"),
			userDB:     os.Getenv("DB_USER"),
			passwordDB: os.Getenv("DB_PASSWORD"),
		},
	}

	conf.Storage.DSN = conf.getDSN()

	return conf, nil
}

func (c *Config) getDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Storage.hostDB,
		c.Storage.userDB,
		c.Storage.passwordDB,
		c.Storage.nameDB,
		c.Storage.portDB,
	)
}

// Storage ..
type Storage struct {
	hostDB     string
	portDB     string
	userDB     string
	passwordDB string
	nameDB     string
	DSN        string
}

type botAntiSpam struct {
	Settings        settings
	WhiteListTags   map[string]struct{}
	WhiteListAuthor []int64
}

type settings struct {
	Token                 string
	AdmToken              string
	OffsetMessageStart    int
	TimeOut               int // for long request to interrapt
	Reties                int // when retries to del message
	TimeOutBetweenRetries time.Duration
}
