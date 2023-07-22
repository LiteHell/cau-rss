package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
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
		if err != err {
			return err
		}

		if len(feed.Items) == 0 {
			return fmt.Errorf("Empty feed")
		}

	case "atom":
		parser := atom.Parser{}
		feed, err := parser.Parse(resp.Body)
		if err != err {
			return err
		}

		if len(feed.Entries) == 0 {
			return fmt.Errorf("Empty feed")
		}

	case "json":
		parser := json.Parser{}
		feed, err := parser.Parse(resp.Body)
		if err != err {
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
		for _, sitename := range []string{"cse", "swedu", "abeek"} {
			url := "http://127.0.0.1:8080/cau/" + sitename + "/" + feedType

			t.Logf("Testing feed: %s", url)
			err := testFeed(url, feedType)
			if err != nil {
				t.Error(err)
			}
		}

		for _, buildingType := range []string{"bluemir", "future_house", "global_house", "all"} {
			url := "http://127.0.0.1:8080/cau/dormitory/seoul/" + buildingType + "/" + feedType

			t.Logf("Testing feed: %s", url)
			err := testFeed(url, feedType)
			if err != nil {
				t.Error(err)
			}
		}

		davinciDormitoryUrl := "http://127.0.0.1:8080/cau/dormitory/davinci/" + feedType
		t.Logf("Testing feed: %s", davinciDormitoryUrl)
		err := testFeed(davinciDormitoryUrl, feedType)
		if err != nil {
			t.Error(err)
		}
	}
}
