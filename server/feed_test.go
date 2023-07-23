package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
	"litehell.info/cau-rss/server"
)

func testFeed(url string, feedType string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	switch feedType {
	case "rss":
		parser := rss.Parser{}
		feed, err := parser.Parse(resp.Body)
		if err != nil {
			return err
		}

		if len(feed.Items) == 0 {
			return fmt.Errorf("Empty feed")
		}

	case "atom":
		parser := atom.Parser{}
		feed, err := parser.Parse(resp.Body)
		if err != nil {
			return err
		}

		if len(feed.Entries) == 0 {
			return fmt.Errorf("Empty feed")
		}

	case "json":
		parser := json.Parser{}
		feed, err := parser.Parse(resp.Body)
		if err != nil {
			return err
		}

		if len(feed.Items) == 0 {
			return fmt.Errorf("Empty feed")
		}
	}

	return nil
}

func TestFeeds(t *testing.T) {
	runServer()

	for _, feedType := range []string{"rss", "atom", "json"} {
		server.LoopForAllSites(func(cw *server.CauWebsite) {
			endpoint := cw.Key + "/" + feedType
			feedType := feedType
			t.Run(endpoint, func(t *testing.T) {
				url := "http://127.0.0.1:8080/cau/" + endpoint

				t.Logf("Testing feed: %s", url)

				t.Parallel()
				err := testFeed(url, feedType)
				if err != nil {
					t.Error(err)
				}
			})
		})
	}
}
