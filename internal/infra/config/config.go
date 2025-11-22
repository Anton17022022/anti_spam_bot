package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config is config app
type Config struct {
	BotAntiSpam botAntiSpam
}

type botAntiSpam struct {
	Settings      settings
	WhiteListTags map[string]struct{}
}

type settings struct {
	Token                 string
	OffsetMessageStart    int
	TimeOut               int // for long request to interrapt
	Reties                int // when retries to del message
	TimeOutBetweenRetries time.Duration
}

// NewConfig return config app instance
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("erorr loading .env file, using default config")
	}

	wlTags := map[string]struct{}{
		"@prolann": {},
		"@Prolann": {},
	}

	return &Config{
		BotAntiSpam: botAntiSpam{
			Settings: settings{
				Token:                 os.Getenv("antispam_bot_token"),
				OffsetMessageStart:    0,
				TimeOut:               60,
				Reties:                3,
				TimeOutBetweenRetries: 10 * time.Second,
			},
			WhiteListTags: wlTags,
		}}, nil
}
