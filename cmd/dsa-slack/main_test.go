package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/renderorange/debian-security-announce/internal/state"
)

func TestMainFlow(t *testing.T) {
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer slackServer.Close()

	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile("../../internal/rss/testdata/dsa-long.xml")
		w.Write(data)
	}))
	defer rssServer.Close()

	dir := t.TempDir()
	stateFile := filepath.Join(dir, "state.json")

	err := run(slackServer.URL, stateFile, rssServer.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, err := loadState(stateFile)
	if err != nil {
		t.Fatalf("failed to load state: %v", err)
	}
	if s.LastDSANumber != 6508 {
		t.Errorf("got last DSA number %d, want 6508", s.LastDSANumber)
	}
}

func TestMainFlow_NoNewItems(t *testing.T) {
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not post to Slack when no new items")
	}))
	defer slackServer.Close()

	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile("../../internal/rss/testdata/dsa-long.xml")
		w.Write(data)
	}))
	defer rssServer.Close()

	dir := t.TempDir()
	stateFile := filepath.Join(dir, "state.json")

	saveState(stateFile, &state.State{LastDSANumber: 6508, LastChecked: "2026-09-20T00:00:00Z"})

	err := run(slackServer.URL, stateFile, rssServer.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
