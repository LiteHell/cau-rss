package server_test

import (
	"os"
	"testing"
	"time"

	"litehell.info/cau-rss/server"
)

func testRedisFor(url string, t *testing.T) {
	baseUrl := "http://127.0.0.1:8080" + url

	t.Logf("Fetching rss first for feed: %s", url)
	err := testFeed(baseUrl+"/rss", "rss")
	if err != nil {
		t.Error(err)
	}

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

func TestRedis(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Logf("Redis test skipped")
		return
	}

	runServer()
	server.InitializeRedis()
	server.ClearCache()

	server.LoopForAllSites(func(cw *server.CauWebsite) {
		key := cw.Key
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			testRedisFor("/cau/"+key, t)
		})
	})
}
