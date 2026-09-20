package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Item struct {
	Title       string
	Link        string
	Date        string
	Description string
	DSANumber   int
}

type rdfFeed struct {
	XMLName xml.Name  `xml:"RDF"`
	Items   []rdfItem `xml:"item"`
}

type rdfItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Date        string `xml:"http://purl.org/dc/elements/1.1/ date"`
	Description string `xml:"description"`
}

var dsaNumberRe = regexp.MustCompile(`DSA-(\d+)-`)

func ParseDSANumber(title string) int {
	matches := dsaNumberRe.FindStringSubmatch(title)
	if len(matches) < 2 {
		return 0
	}
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}
	return num
}

func Fetch(feedURL string) ([]Item, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RSS feed returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSS feed body: %w", err)
	}

	return parseXML(data)
}

func parseXML(data []byte) ([]Item, error) {
	var feed rdfFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	items := make([]Item, 0, len(feed.Items))
	for _, ri := range feed.Items {
		title := strings.TrimSpace(ri.Title)
		items = append(items, Item{
			Title:       title,
			Link:        ri.Link,
			Date:        ri.Date,
			Description: strings.TrimSpace(ri.Description),
			DSANumber:   ParseDSANumber(title),
		})
	}

	return items, nil
}

func FilterNew(items []Item, lastDSANumber int) []Item {
	var result []Item
	for _, item := range items {
		if item.DSANumber > lastDSANumber {
			result = append(result, item)
		}
	}
	return result
}
