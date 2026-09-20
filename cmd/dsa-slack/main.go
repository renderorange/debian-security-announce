package main

import (
	"fmt"
	"log"
	"time"

	"github.com/renderorange/debian-security-announce/internal/config"
	"github.com/renderorange/debian-security-announce/internal/rss"
	"github.com/renderorange/debian-security-announce/internal/slack"
	"github.com/renderorange/debian-security-announce/internal/state"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := run(cfg.WebhookURL, cfg.StateFile, cfg.RSSFeedURL); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(webhookURL, stateFile, rssFeedURL string) error {
	s, err := loadState(stateFile)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	items, err := rss.Fetch(rssFeedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch RSS feed: %w", err)
	}

	newItems := rss.FilterNew(items, s.LastDSANumber)
	if len(newItems) == 0 {
		return nil
	}

	for _, item := range newItems {
		if err := slack.Post(webhookURL, item); err != nil {
			return fmt.Errorf("failed to post DSA %d to Slack: %w", item.DSANumber, err)
		}
	}

	maxDSA := 0
	for _, item := range newItems {
		if item.DSANumber > maxDSA {
			maxDSA = item.DSANumber
		}
	}
	s.LastDSANumber = maxDSA
	s.LastChecked = time.Now().UTC().Format(time.RFC3339)

	if err := saveState(stateFile, s); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

func loadState(filePath string) (*state.State, error) {
	return state.Load(filePath)
}

func saveState(filePath string, s *state.State) error {
	return state.Save(filePath, s)
}
