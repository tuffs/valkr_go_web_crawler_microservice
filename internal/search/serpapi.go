// internal/search/serpapi.go
package search

import (
	"context"
	"strconv"
)

type Page struct {
	URL          string
	Title        string
	Description  string
	ThumbnailURL string
}

func Search(ctx context.Context, apiKey, query string, num int) ([]Page, error) {
	params := map[string]string{
		"q":        query,
		"num":      strconv.Itoa(num),
		"api_key":  apiKey,
		"engine":   "google",
		"gl":       "us",
		"hl":       "en",
	}

	search := google_search_results.NewGoogleSearch(params, "")
	result, err := search.GetJSON(ctx)
	if err != nil {
		return nil, err
	}

	var pages []Page
	for _, r := range result.GetArray("organic_results") {
		m := r.(map[string]interface{})
		link := getString(m, "link")
		title := getString(m, "title")
		snippet := getString(m, "snippet")
		thumb := ""

		// Try to get thumbnail
		if thumbObj, ok := m["thumbnail"].(map[string]interface{}); ok {
			thumb = getString(thumbObj, "src")
		}

		if link != "" && title != "" {
			pages = append(pages, Page{
				URL:          link,
				Title:        truncate(title, 255),
				Description:  truncate(snippet, 255),
				ThumbnailURL: thumb,
			})
		}
	}

	return pages, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}