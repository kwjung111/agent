package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	scrapeInterval time.Duration
}

func (c *Config) initConfig() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to get .env file : %v", err)
	}

	intervalStr := os.Getenv("SCRAPE_INTERVAL")

	duration, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("Error Parsing scrape_interval value : %v", err)
	}
	c.scrapeInterval = duration
}

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{}
		config.initConfig()
	})
	return config
}

func (c *Config) GetScrapeInterval() time.Duration {
	return config.scrapeInterval
}
