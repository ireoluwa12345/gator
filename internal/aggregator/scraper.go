package aggregator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ireoluwa12345/gator/internal/database"
	"github.com/ireoluwa12345/gator/internal/domain"
	"github.com/ireoluwa12345/gator/internal/rss"
	"github.com/lib/pq"
)

func ScrapeFeeds(s *domain.State) error {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.DB.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)

	for _, item := range feedData.Channel.Item {
		formattedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return fmt.Errorf("couldn't formate published date")
		}
		_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			Title:       item.Title,
			FeedID:      feed.ID,
			Description: sql.NullString{String: item.Description, Valid: true},
			Url:         item.Link,
			PublishedAt: sql.NullTime{Time: formattedTime, Valid: true},
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					continue
				}
			} else {
				fmt.Printf("error occurred: %v\n", err)
			}
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}
