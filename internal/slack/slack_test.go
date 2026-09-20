package slack

import (
	"testing"

	"github.com/renderorange/debian-security-announce/internal/rss"
)

func TestFormatMessage(t *testing.T) {
	item := rss.Item{
		Title:       "DSA-6508-1 chromium - security update",
		Link:        "https://lists.debian.org/debian-security-announce/2026/msg00420.html",
		Date:        "2026-09-19",
		Description: "Security issues were discovered in Chromium.\n\nhttps://security-tracker.debian.org/tracker/DSA-6508-1",
		DSANumber:   6508,
	}

	msg := FormatMessage(item)

	if msg == "" {
		t.Fatal("got empty message")
	}
	if !contains(msg, "DSA-6508-1") {
		t.Error("message should contain DSA number")
	}
	if !contains(msg, "chromium") {
		t.Error("message should contain package name")
	}
	if !contains(msg, item.Link) {
		t.Error("message should contain archive link")
	}
	if !contains(msg, "Tracker") {
		t.Error("message should contain tracker link text")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
