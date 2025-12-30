package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ireoluwa12345/gator/internal/database"
	"github.com/ireoluwa12345/gator/internal/domain"
)

func HandleAddFeed(s *domain.State, cmd Command, user database.GetUserRow) error {
	if len(cmd.CommandArgs) < 2 {
		return fmt.Errorf("the handler expects the name and url")
	}

	name := cmd.CommandArgs[0]
	url := cmd.CommandArgs[1]

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return fmt.Errorf("error occurred while creating the feed: %w", err)
	}

	_, err = s.DB.GetFollowedFeedByID(context.Background(), feed.ID)

	if err != nil {
		return fmt.Errorf("error occurred while fetching followed feed: %w", err)
	} else {
		_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UserID: feed.UserID, FeedID: feed.ID})

		if err != nil {
			return fmt.Errorf("error occurred while following the feed: %w", err)
		}
	}

	fmt.Printf("%s added to feeds", feed.Name)

	return nil
}

func HandleFeeds(s *domain.State, cmd Command) error {
	feeds, err := s.DB.GetAllFeeds(context.Background())

	if err != nil {
		return fmt.Errorf("error occurred while fetching feeds: %w", err)
	}

	for _, feed := range feeds {
		if feed.UserName.Valid {
			fmt.Printf("- %s by %s (%s)  \n", feed.Name, feed.UserName.String, feed.Url)
			continue
		}
		fmt.Printf("- %s (%s)  \n", feed.Name, feed.Url)
	}

	return nil
}

func HandleFollow(s *domain.State, cmd Command, user database.GetUserRow) error {
	if len(cmd.CommandArgs) < 1 {
		return fmt.Errorf("the handler expects the feed name")
	}

	feedName := cmd.CommandArgs[0]

	feed, err := s.DB.GetFeedByName(context.Background(), feedName)
	if err != nil {
		return fmt.Errorf("error occurred while fetching feed: %w", err)
	}

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error occurred while following the feed: %w", err)
	}

	fmt.Printf("%s now following %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func HandleFollowing(s *domain.State, cmd Command, user database.GetUserRow) error {
	followingFeeds, err := s.DB.GetFeedsFollowForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("error occurred while fetching following feeds: %w", err)
	}

	for _, followingFeed := range followingFeeds {
		fmt.Printf("%s", followingFeed.FeedName)
	}

	return nil
}

func HandleUnfollow(s *domain.State, cmd Command, user database.GetUserRow) error {
	if len(cmd.CommandArgs) < 1 {
		return fmt.Errorf("the handler expects the feed url")
	}

	feedUrl := cmd.CommandArgs[0]

	feed, err := s.DB.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error occurred while fetching feed: %w", err)
	}

	err = s.DB.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error occurred while unfollowing the feed: %w", err)
	}

	fmt.Printf("%s has unfollowed %s\n", user.Name, feed.Name)

	return nil
}

func HandleBrowse(s *domain.State, cmd Command, user database.GetUserRow) error {
	posts, err := s.DB.FetchUserPosts(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("error occured: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nUrl: %s\n", post.Title, post.Url)
	}

	return nil
}
