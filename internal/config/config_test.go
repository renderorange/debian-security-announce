package config

import (
	"os"
	"testing"
)

func TestLoad_RequiredWebhookURL(t *testing.T) {
	os.Unsetenv("DSA_SLACK_WEBHOOK_URL")
	_, err := Load()
	if err == nil {
		t.Error("expected error when DSA_SLACK_WEBHOOK_URL not set")
	}
}

func TestLoad_WithWebhookURL(t *testing.T) {
	os.Setenv("DSA_SLACK_WEBHOOK_URL", "https://hooks.slack.com/test")
	defer os.Unsetenv("DSA_SLACK_WEBHOOK_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WebhookURL != "https://hooks.slack.com/test" {
		t.Errorf("got webhook URL %q, want %q", cfg.WebhookURL, "https://hooks.slack.com/test")
	}
}

func TestLoad_DefaultStateFile(t *testing.T) {
	os.Setenv("DSA_SLACK_WEBHOOK_URL", "https://hooks.slack.com/test")
	defer os.Unsetenv("DSA_SLACK_WEBHOOK_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile == "" {
		t.Error("expected default state file path, got empty string")
	}
}

func TestLoad_CustomStateFile(t *testing.T) {
	os.Setenv("DSA_SLACK_WEBHOOK_URL", "https://hooks.slack.com/test")
	os.Setenv("DSA_STATE_FILE", "/tmp/custom-state.json")
	defer os.Unsetenv("DSA_SLACK_WEBHOOK_URL")
	defer os.Unsetenv("DSA_STATE_FILE")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile != "/tmp/custom-state.json" {
		t.Errorf("got state file %q, want %q", cfg.StateFile, "/tmp/custom-state.json")
	}
}

func TestLoad_DefaultRSSFeedURL(t *testing.T) {
	os.Setenv("DSA_SLACK_WEBHOOK_URL", "https://hooks.slack.com/test")
	defer os.Unsetenv("DSA_SLACK_WEBHOOK_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RSSFeedURL != "https://www.debian.org/security/dsa-long" {
		t.Errorf("got RSS feed URL %q, want %q", cfg.RSSFeedURL, "https://www.debian.org/security/dsa-long")
	}
}
