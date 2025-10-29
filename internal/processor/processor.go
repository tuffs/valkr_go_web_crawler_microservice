package processor

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/db"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/search"
)

func ProcessTrend(ctx context.Context, trend string, page []search.Page) {
	// 1. Insert into crawl_queues
	urls := make([]string, len(pages))
	for i, p := range pages {
		urls[i] = p.URL
	}
	if err := db.InsertCrawlQueue(ctx, trend, urls); err != nil {
		return err
	}

	// 2. Upsert each webpage
	for _, p := range pages {
		if err := db.UpsertWebpage(ctx, p); err != nil {
			log.Err(err).Str("url", p.URL).Msg("failed to upsert")
			// continue
		}
	}

	log.Info().Str("trend", trend).Int("pages", len(pages)).Msg("processed")
	return nil
}