package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"litehell.info/cau-rss/cau_parser"
)

func serveFeed(ctx *gin.Context) {
	feed_tmp, hasFeeds := ctx.Get("feed")
	articles_tmp, hasArticles := ctx.Get("articles")

	if !hasFeeds || !hasArticles {
		return
	}

	feed := feed_tmp.(*feeds.Feed)
	articles := articles_tmp.([]cau_parser.CAUArticle)

	var feedStr string
	var feedErr error
	var contentType string

	switch ctx.Param("feedType") {
	case "rss":
		feedStr, feedErr = generateFeed(feed, articles, RSS)
		contentType = "application/rss+xml"
	case "atom":
		feedStr, feedErr = generateFeed(feed, articles, ATOM)
		contentType = "application/atom+xml"
	case "json":
		feedStr, feedErr = generateFeed(feed, articles, JSON)
		contentType = "application/feed+json"
	default:
		fmt.Fprintf(os.Stderr, "unsupported feed type: %s", ctx.Param("feedType"))
		ctx.HTML(404, "404.html", gin.H{})
		return
	}

	if feedErr != nil {
		fmt.Fprint(os.Stderr, feedErr)
		ctx.HTML(500, "feed-gen-error.html", gin.H{})
		return
	}

	saveCache(feedStr, contentType, ctx)

	ctx.Writer.Header().Set("Content-Type", contentType)
	ctx.String(200, feedStr)
}
