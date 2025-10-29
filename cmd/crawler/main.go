package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/config"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/db"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/processor"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/search"
	"github.com/tuffs/valkr_go_web_crawler_microservice/internal/trends"
	"github.com/urfave/cli/v2"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zeroLog.ConsoleWriter{Out: os.Stderr})

	app := &cli.App{
		Name: 'valkr-crawler',
		Usage: "Fetch Google Trends -> Top 100 Pages -> Store in Neon Postgresql",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: 'date',
				Usage: 'Date in YYYY-MM-DD (default: today)',
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(msg).Msg('crawler failed')
	}
}

func run(c *cli.Context) error {
	cfg := config.Load()

	// Override data from CLI
	if dateStr := c.String('date'); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
		cfg.TargetDate = parsed
	}

	if cfg.DatabaseURL == "" {
		return cli.Exit("DATABASE_URL is required", 1)
	}
	if cfg.SerpAPIKEY == "" {
		return cli.Exit("SERPAPI_KEY is required", 1)
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// 1. Get top 10 trends
	trends, err := trends.GetDailyTopTrends(ctx, cfg.TrendsGeo, cfg.TargetDate)
	if err != nil {
		return err
	}
	if len(trends) == 0 {
		return err
	}
	if len(trends) == 0 {
		log.Warn().Msg("no trends found for date")
		return nil
	}

	log.Info().Str("date", cfg.TargetDate.Format("2006-01-02")).Int("trends", len(trends)).Msg("fetched trends")

	// 2. Process each trend
	for _, trend := range trends {
		select {
		case <-ctx.Done():
				return ctx.Err()
		default:
		}

		log.Info().Str("trend", trend).Msg("searching")
		pages, err := search.Search(ctx, cfg.SerpAPIKey, trend, cfg.ResultsPerTrend)
		if err != nil {
			log.Err(err).Str("trend", trend).Msg("search failed")
			continue
		}

		if err := processor.ProcessTrend(ctx, trend, pages); err != nil {
			log.Err(err).Str("trend", trend).Msg("processing failed")
		}
	}

	return nil
}