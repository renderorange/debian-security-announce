package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	WebhookURL string
	StateFile  string
	RSSFeedURL string
}

func Load() (*Config, error) {
	webhookURL := os.Getenv("DSA_SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return nil, fmt.Errorf("DSA_SLACK_WEBHOOK_URL environment variable is required")
	}

	stateFile := os.Getenv("DSA_STATE_FILE")
	if stateFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		stateFile = filepath.Join(home, ".config", "debian-security-announce", "state.json")
	}

	rssFeedURL := os.Getenv("DSA_RSS_FEED_URL")
	if rssFeedURL == "" {
		rssFeedURL = "https://www.debian.org/security/dsa-long"
	}

	return &Config{
		WebhookURL: webhookURL,
		StateFile:  stateFile,
		RSSFeedURL: rssFeedURL,
	}, nil
}
