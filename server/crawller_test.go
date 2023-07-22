package server_test

import (
	"os"
	"testing"
	"time"

	"litehell.info/cau-rss/server"
)

func testCrawllerFor(url string, t *testing.T) {
	baseUrl := "http://127.0.0.1:8080" + url

	for _, i := range []string{"rss", "atom", "json"} {
		t.Logf("Testing redis feed: %s", url)
		timeBeforeTest := time.Now()
		err := testFeed(baseUrl+"/"+i, i)
		elapsedTime := time.Now().Sub(timeBeforeTest) / time.Millisecond
		if err != nil {
			t.Error(err)
		} else if elapsedTime > time.Millisecond*100 {
			t.Errorf("Too slow for feed: %s (%s), %s", url, i, elapsedTime.String())
		}
	}
}

func TestCrawller(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Logf("Redis test skipped")
		return
	}

	server.InitializeRedis()
	server.ClearCache()
	server.StartCrawller()
	runServer()

	time.Sleep(5)
	for _, sitename := range []string{"cse", "ai", "swedu", "abeek"} {
		testCrawllerFor("/cau/"+sitename, t)
	}

	for _, buildingType := range []string{"bluemir", "future_house", "global_house", "all"} {
		testCrawllerFor("/cau/dormitory/seoul/"+buildingType, t)
	}
	testCrawllerFor("/cau/dormitory/davinci", t)
}
