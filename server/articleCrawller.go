package server

import (
	"fmt"
	"os"
	"time"
)

func StartCrawller() {
	if isRedisAvailable() {
		var crawlChannel chan string = make(chan string)
		go commandSender(crawlChannel)

		for i := 0; i < 4; i++ {
			go fetchWorker(crawlChannel)
		}
	}
}

func fetchWorker(ch <-chan string) {
	for {
		for i := range ch {
			fmt.Printf("Crawlling %s\n", i)
			articles, articlesErr := fetchArticlesForKey(i)
			if articlesErr != nil {
				fmt.Fprintf(os.Stderr, "error while crawling(%s): %s\n", i, articlesErr)
				continue
			}

			err := saveArticlesCahe(i, articles)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error while saving into cache(%s): %s\n", i, err)
			} else {
				fmt.Printf("Done crawling: %s\n", i)
			}
		}
	}
}

func commandSender(ch chan<- string) {
	for {
		LoopForAllSites(func(cw *CauWebsite) {
			if cw.Key == "dormitory/seoul/all" {
				return
			}
			ch <- cw.Key
		})
		time.Sleep(time.Minute * 3)
	}
}
