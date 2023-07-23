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
	if os.Getenv("REDIS_ENABLED") != "true" {
		t.Skip("Redis not available, skipping...")
		return
	}

	server.InitializeRedis()
	server.ClearCache()
	server.StartCrawller()
	runServer()

	time.Sleep(5)
	server.LoopForAllSites(func(cw *server.CauWebsite) {
		testCrawllerFor("/cau/"+cw.Key, t)
	})
}
