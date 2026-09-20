package rss

import (
	"os"
	"testing"
)

func TestParseDSANumber(t *testing.T) {
	tests := []struct {
		title string
		want  int
	}{
		{"DSA-6508-1 chromium - security update", 6508},
		{"DSA-6507-1 unbound - security update", 6507},
		{"DSA-6496-2 nginx - regression update", 6496},
		{"invalid title", 0},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := ParseDSANumber(tt.title)
			if got != tt.want {
				t.Errorf("ParseDSANumber(%q) = %d, want %d", tt.title, got, tt.want)
			}
		})
	}
}

func TestFetch_FromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/dsa-long.xml")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	items, err := parseXML(data)
	if err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}

	if items[0].Title != "DSA-6508-1 chromium - security update" {
		t.Errorf("got title %q, want %q", items[0].Title, "DSA-6508-1 chromium - security update")
	}
	if items[0].DSANumber != 6508 {
		t.Errorf("got DSA number %d, want 6508", items[0].DSANumber)
	}
	if items[0].Link != "https://lists.debian.org/debian-security-announce/2026/msg00420.html" {
		t.Errorf("got link %q, want %q", items[0].Link, "https://lists.debian.org/debian-security-announce/2026/msg00420.html")
	}
}

func TestFilterNew(t *testing.T) {
	items := []Item{
		{DSANumber: 6508, Title: "DSA-6508-1"},
		{DSANumber: 6507, Title: "DSA-6507-1"},
		{DSANumber: 6506, Title: "DSA-6506-1"},
	}

	filtered := FilterNew(items, 6507)
	if len(filtered) != 1 {
		t.Fatalf("got %d items, want 1", len(filtered))
	}
	if filtered[0].DSANumber != 6508 {
		t.Errorf("got DSA number %d, want 6508", filtered[0].DSANumber)
	}
}

func TestFilterNew_AllNew(t *testing.T) {
	items := []Item{
		{DSANumber: 6508, Title: "DSA-6508-1"},
		{DSANumber: 6507, Title: "DSA-6507-1"},
	}

	filtered := FilterNew(items, 0)
	if len(filtered) != 2 {
		t.Fatalf("got %d items, want 2", len(filtered))
	}
}

func TestFilterNew_NoneNew(t *testing.T) {
	items := []Item{
		{DSANumber: 6508, Title: "DSA-6508-1"},
	}

	filtered := FilterNew(items, 6508)
	if len(filtered) != 0 {
		t.Fatalf("got %d items, want 0", len(filtered))
	}
}
