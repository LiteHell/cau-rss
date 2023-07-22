package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"litehell.info/cau-rss/cau_parser"
)

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

	switch key {
	case "cse":
		articles, articlesErr = cau_parser.ParseCSE()
	case "swedu":
		articles, articlesErr = cau_parser.ParseSWEDU()
	case "abeek":
		articles, articlesErr = cau_parser.ParseABEEK()
	case "dormitory/davinci":
		articles, articlesErr = cau_parser.ParseDormitory(cau_parser.DORMITORY_DAVINCI)
	case "dormitory/seoul/bluemir":
		articles, articlesErr = cau_parser.ParseDormitory(cau_parser.DORMITORY_BLUEMIR)
	case "dormitory/seoul/future_house":
		articles, articlesErr = cau_parser.ParseDormitory(cau_parser.DORMITORY_FUTURE_HOUSE)
	case "dormitory/seoul/global_house":
		articles, articlesErr = cau_parser.ParseDormitory(cau_parser.DORMITORY_GLOBAL_HOUSE)
	default:
		panic(fmt.Errorf("Unknown articles key: %s", key))
	}

	if articlesErr != nil {
		return nil, articlesErr
	}

	if isRedisAvailable() {
		serialized, err := json.Marshal(articles)
		if err != nil {
			return articles, err
		}
		redisClient.Set(redisCtx, "articles/"+key, string(serialized), time.Minute*5)
	}

	return articles, nil
}
