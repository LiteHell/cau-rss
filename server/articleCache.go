package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"litehell.info/cau-rss/cau_parser"
)

func fetchArticlesForKey(key string) ([]cau_parser.CAUArticle, error) {
	switch key {
	case "cse":
		return cau_parser.ParseCSE()
	case "swedu":
		return cau_parser.ParseSWEDU()
	case "abeek":
		return cau_parser.ParseABEEK()
	case "dormitory/davinci":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_DAVINCI)
	case "dormitory/seoul/bluemir":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_BLUEMIR)
	case "dormitory/seoul/future_house":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_FUTURE_HOUSE)
	case "dormitory/seoul/global_house":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_GLOBAL_HOUSE)
	default:
		panic(fmt.Errorf("Unknown articles key: %s", key))
	}
}

func saveArticlesCahe(key string, articles []cau_parser.CAUArticle) error {
	if isRedisAvailable() {
		serialized, err := json.Marshal(articles)
		if err != nil {
			return err
		}
		return redisClient.Set(redisCtx, "articles/"+key, string(serialized), time.Minute*5).Err()
	}
	return nil
}

func getArticlesWithCache(key string) ([]cau_parser.CAUArticle, error) {
	var articles []cau_parser.CAUArticle
	var articlesErr error

	if isRedisAvailable() {
		cache, err := redisClient.Get(redisCtx, "articles/"+key).Result()

		if err != nil && err != redis.Nil {
			return nil, err
		} else if err == nil {
			json.Unmarshal([]byte(cache), &articles)
			return articles, nil
		}
	}

	articles, articlesErr = fetchArticlesForKey(key)

	if articlesErr != nil {
		return nil, articlesErr
	}

	err := saveArticlesCahe(key, articles)

	return articles, err
}
