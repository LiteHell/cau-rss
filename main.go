package main

import (
	"os"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"litehell.info/cau-rss/cau_parser"
	"litehell.info/cau-rss/server"
)

type CrawlSuccessItem struct {
	Articles  []cau_parser.CAUArticle `json:"articles"`
	SiteInfo  server.CauWebsite       `json:"site_info"`
	Timestamp int64                   `json:"timestamp"`
}

type CrwalFailureItem struct {
	SiteInfo  server.CauWebsite
	Timestamp int64 `json:"timestamp"`
}

type FeedDataResponse struct {
	Success []CrawlSuccessItem `json:"success"`
	Failure []CrwalFailureItem `json:"failure"`
}

func getAllSeoulDormitoryArticles(siteInfo *server.CauWebsite, items *[]CrawlSuccessItem) []cau_parser.CAUArticle {
	articles := make([]cau_parser.CAUArticle, 0)

	for _, item := range *items {
		if item.SiteInfo.Key == "dormitory/seoul/bluemir" ||
			item.SiteInfo.Key == "dormitory/seoul/future_house" ||
			item.SiteInfo.Key == "dormitory/seoul/global_house" {
			articles = append(articles, item.Articles...)
		}
	}

	sort.Slice(articles, func(a, b int) bool {
		return articles[a].Date.Before(articles[b].Date)
	})

	return articles
}

func HandleLambdaEvent() (*FeedDataResponse, error) {
	success := make([]CrawlSuccessItem, 0)
	failures := make([]CrwalFailureItem, 0)
	var allSeoulDormitory server.CauWebsite

	server.LoopForAllSites(func(site *server.CauWebsite) {
		if site.Key == "dormitory/seoul/all" {
			allSeoulDormitory = *site
			return
		}
		articles, articlesErr := server.FetchArticlesForKey(site.Key)
		if articlesErr != nil {
			failures = append(failures, CrwalFailureItem{
				SiteInfo:  *site,
				Timestamp: time.Now().Unix(),
			})
			return
		} else {
			success = append(success, CrawlSuccessItem{
				SiteInfo:  *site,
				Articles:  articles,
				Timestamp: time.Now().Unix(),
			})
		}
	})

	success = append(success, CrawlSuccessItem{
		SiteInfo:  allSeoulDormitory,
		Articles:  getAllSeoulDormitoryArticles(&allSeoulDormitory, &success),
		Timestamp: time.Now().Unix(),
	})

	distDir, err := os.MkdirTemp("", "*******")
	if err != nil {
		panic(err)
	}

	res := FeedDataResponse{
		Success: success,
		Failure: failures,
	}

	generateStatic(distDir, &res)
	uploadS3(distDir)

	return &res, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
