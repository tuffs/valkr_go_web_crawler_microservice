package trends

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

func GetDailyTopTrends(ctx context.Context, geo string, date time.Time) ([]string, error) {
	// Format: https://trends.google.com/trends/trendingsearches/daily?geo=US&date=2025-10-27
	url := fmt.Sprintf(
		"https://trends.google.com/trends/trendingsearches/daily?geo=%s&date=%s",
		geo, date.Format("2006-01-02"),
	)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 15 + time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var trends []string
	doc.Find(".title-container .title").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := strings.TrimSpace(s.Text())
		if title != "" && len(trends) < 10 {
			trends = append(trends, title)
		}
		return len(trends) < 10
	})

	if len(trends) == 0 {
		log.Warn().Str("url", url).Msg("no trends found - page structure may have changed")
	}

	return trends, nil
}