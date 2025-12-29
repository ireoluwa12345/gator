package cli

import (
	"fmt"
	"time"

	"github.com/ireoluwa12345/gator/internal/aggregator"
	"github.com/ireoluwa12345/gator/internal/domain"
)

func HandleAgg(s *domain.State, cmd Command) error {
	if len(cmd.CommandArgs) < 1 {
		return fmt.Errorf("the handler requires a time_between_reqs")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.CommandArgs[0])

	if err != nil {
		return fmt.Errorf("error occurred while parsing time_between_reqs: %w", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	fmt.Printf("Collecting feeds every %v\n\n", timeBetweenReqs)
	for ; ; <-ticker.C {
		err := aggregator.ScrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
			break
		}
	}

	return nil
}
