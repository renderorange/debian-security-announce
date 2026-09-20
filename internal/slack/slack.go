package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/renderorange/debian-security-announce/internal/rss"
)

func FormatMessage(item rss.Item) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("*%s*\n", item.Title))
	b.WriteString(fmt.Sprintf("_security update_\n\n"))
	b.WriteString(item.Description)
	b.WriteString("\n\n")

	trackerURL := extractTrackerURL(item.Description)
	if trackerURL != "" {
		b.WriteString(fmt.Sprintf("<%s|Tracker> | <%s|Archive>", trackerURL, item.Link))
	} else {
		b.WriteString(fmt.Sprintf("<%s|Archive>", item.Link))
	}

	return b.String()
}

func extractTrackerURL(description string) string {
	lines := strings.Split(description, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "security-tracker.debian.org/tracker/") {
			parts := strings.Fields(line)
			for _, p := range parts {
				if strings.HasPrefix(p, "https://security-tracker.debian.org/tracker/") {
					return p
				}
			}
		}
	}
	return ""
}

func Post(webhookURL string, item rss.Item) error {
	msg := FormatMessage(item)

	payload := map[string]string{
		"text": msg,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to post to Slack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}
