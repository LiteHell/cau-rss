package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"litehell.info/cau-rss/cau_parser"
)

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("html/*.html")

	server.GET("/cau/notice", func(ctx *gin.Context) {
		ctx.Redirect(308, "https://www.cau.ac.kr/cms/FR_PRO_CON/BoardRss.do?pageNo=1&pagePerCnt=15&MENU_ID=100&SITE_NO=2&BOARD_SEQ=4&S_CATE_SEQ=&BOARD_TYPE=C0301&BOARD_CATEGORY_NO=&P_TAB_NO=&TAB_NO=&P_CATE_SEQ=&CATE_SEQ=&SEARCH_FLD=SUBJECT&SEARCH=")
	})
	server.GET("/cau/:siteType/:feedType", func(ctx *gin.Context) {
		var feed *feeds.Feed
		var articles []cau_parser.CAUArticle
		var articlesErr error
		switch ctx.Param("siteType") {
		case "sw":
			ctx.Redirect(308, "/cau/swedu/"+ctx.Param("feedType"))
			return
		case "cse":
			feed = &feeds.Feed{
				Title:       "중앙대학교 소프트웨어학부 공지사항",
				Link:        &feeds.Link{Href: "https://cse.cau.ac.kr"},
				Description: "중앙대학교 소프트웨어학부의 공지사항입니다",
				Created:     time.Now(),
			}
			articles, articlesErr = cau_parser.ParseCSE()
		case "swedu":
			feed = &feeds.Feed{
				Title:       "중앙대학교 다빈치SW교욱원 공지사항",
				Link:        &feeds.Link{Href: "https://swedu.cau.ac.kr"},
				Description: "중앙대학교 다빈치SW교욱원의 공지사항입니다",
				Created:     time.Now(),
			}
			articles, articlesErr = cau_parser.ParseSWEDU()
		case "abeek":
			feed = &feeds.Feed{
				Title:       "중앙대학교 공학교육혁신센터 공지사항",
				Link:        &feeds.Link{Href: "https://abek.cau.ac.kr"},
				Description: "중앙대학교 공학교육혁신센터의 공지사항입니다",
				Created:     time.Now(),
			}
			articles, articlesErr = cau_parser.ParseABEEK()
		default:
			fmt.Fprintf(os.Stderr, "unsupported website: %s", ctx.Param("siteType"))
			ctx.HTML(404, "404.html", gin.H{})
			return
		}
		if articlesErr != nil {
			fmt.Fprint(os.Stderr, articlesErr)
			ctx.HTML(500, "article-fetch-error.html", gin.H{})
			return
		}

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

		ctx.Writer.Header().Set("Content-Type", contentType)
		ctx.String(200, feedStr)
	})

	server.Run()
}
