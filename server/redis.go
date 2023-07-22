package server

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client = nil
var redisCtx context.Context

const redis_feed_cache_key_id = "redis_cache_key"

func InitializeRedis() {
	if os.Getenv("REDIS_ENABLED") != "true" {
		return
	}

	redisCtx = context.Background()
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 0
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       db,                          // use default DB
	})
}

func getRedisFeedCacheKey(ctx *gin.Context) (string, string, bool) {
	cacheKey_tmp, hasCacheKey := ctx.Get(redis_feed_cache_key_id)
	if !hasCacheKey {
		return "", "", false
	}

	cacheKey := cacheKey_tmp.(string)
	cacheMimeKey := cacheKey + "_mime"

	return cacheKey, cacheMimeKey, true
}

func setRedisFeedCacheKey(ctx *gin.Context, key string) {
	ctx.Set(redis_feed_cache_key_id, key)
}

func saveFeedCache(feed string, mime string, ctx *gin.Context) {
	if redisClient == nil {
		return
	}

	cacheKey, cacheMimeKey, hasCacheKey := getRedisFeedCacheKey(ctx)
	if !hasCacheKey {
		return
	}

	redisClient.Set(redisCtx, cacheKey, feed, time.Minute*5)
	redisClient.Set(redisCtx, cacheMimeKey, mime, time.Minute*5)
}

func serveCachedFeed(ctx *gin.Context) {
	if redisClient == nil {
		return
	}

	cacheKey, cacheMimeKey, hasCacheKey := getRedisFeedCacheKey(ctx)
	if !hasCacheKey {
		return
	}

	cache, err := redisClient.Get(redisCtx, cacheKey).Result()
	cacheMime, mimeErr := redisClient.Get(redisCtx, cacheMimeKey).Result()

	if err != nil && err != redis.Nil {
		panic(err)
	} else if mimeErr != nil && mimeErr != redis.Nil {
		panic(mimeErr)
	} else if err == redis.Nil || mimeErr == redis.Nil {
		return
	}

	ctx.Writer.Header().Add("Content-Type", cacheMime)
	ctx.String(200, cache)
	ctx.Abort()
}
