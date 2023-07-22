package server

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"litehell.info/cau-rss/cau_parser"
)

func getAllSeoulDormitoryArticles() ([]cau_parser.CAUArticle, error) {
	articles := []cau_parser.CAUArticle{}
	var articlesErr error
	for _, i := range []string{"bluemir", "future_house", "global_house"} {
		articlesNow, articlesErrNow := getArticlesWithCache("dormitory/seoul/" + i)
		if articlesErrNow != nil {
			articlesErr = articlesErrNow
		}
		articles = append(articles, articlesNow...)
	}

	if articlesErr != nil {
		return nil, articlesErr
	} else {
		return articles, nil
	}
}

func CreateServer() *gin.Engine {
	server := gin.Default()
	server.LoadHTMLGlob("html/*.html")

	server.GET("/index.html", func(ctx *gin.Context) {
		ctx.Redirect(308, "/")
	})
	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

	server.GET("/cau/notice", func(ctx *gin.Context) {
		ctx.Redirect(308, "https://www.cau.ac.kr/cms/FR_PRO_CON/BoardRss.do?pageNo=1&pagePerCnt=15&MENU_ID=100&SITE_NO=2&BOARD_SEQ=4&S_CATE_SEQ=&BOARD_TYPE=C0301&BOARD_CATEGORY_NO=&P_TAB_NO=&TAB_NO=&P_CATE_SEQ=&CATE_SEQ=&SEARCH_FLD=SUBJECT&SEARCH=")
	})
	server.GET("/cau/:siteType/:feedType", func(ctx *gin.Context) {
		setRedisFeedCacheKey(ctx, ctx.Request.URL.Path)
	}, serveCachedFeed, func(ctx *gin.Context) {
		var feed *feeds.Feed
		var articles []cau_parser.CAUArticle = []cau_parser.CAUArticle{}
		var articlesErr error
		siteType := ctx.Param("siteType")
		switch siteType {
		case "sw":
			ctx.Redirect(308, "/cau/swedu/"+ctx.Param("feedType"))
			return
		case "dormitory":
			ctx.Redirect(308, "/cau/dormitory/seoul/all/"+ctx.Param("feedType"))
			return
		case "cse":
			feed = &feeds.Feed{
				Title:       "중앙대학교 소프트웨어학부 공지사항",
				Link:        &feeds.Link{Href: "https://cse.cau.ac.kr"},
				Description: "중앙대학교 소프트웨어학부의 공지사항입니다",
				Created:     time.Now(),
			}
		case "swedu":
			feed = &feeds.Feed{
				Title:       "중앙대학교 다빈치SW교욱원 공지사항",
				Link:        &feeds.Link{Href: "https://swedu.cau.ac.kr"},
				Description: "중앙대학교 다빈치SW교욱원의 공지사항입니다",
				Created:     time.Now(),
			}
		case "abeek":
			feed = &feeds.Feed{
				Title:       "중앙대학교 공학교육혁신센터 공지사항",
				Link:        &feeds.Link{Href: "https://abek.cau.ac.kr"},
				Description: "중앙대학교 공학교육혁신센터의 공지사항입니다",
				Created:     time.Now(),
			}
		default:
			fmt.Fprintf(os.Stderr, "unsupported website: %s", ctx.Param("siteType"))
			ctx.HTML(404, "404.html", gin.H{})
			return
		}

		articles, articlesErr = getArticlesWithCache(siteType)
		if articlesErr != nil {
			fmt.Fprint(os.Stderr, articlesErr)
			ctx.HTML(500, "article-fetch-error.html", gin.H{})
			return
		}

		ctx.Set("feed", feed)
		ctx.Set("articles", articles)
	}, serveFeed)
	server.GET("/cau/dormitory/davinci/:feedType", func(ctx *gin.Context) {
		setRedisFeedCacheKey(ctx, ctx.Request.URL.Path)
	}, serveCachedFeed, func(ctx *gin.Context) {
		articles, articlesErr := getArticlesWithCache("dormitory/davinci")

		if articlesErr != nil {
			fmt.Fprint(os.Stderr, articlesErr)
			ctx.HTML(500, "article-fetch-error.html", gin.H{})
			return
		}

		feed := &feeds.Feed{
			Title:       "중앙대학교 다빈치캠퍼스 기숙사 공지사항",
			Link:        &feeds.Link{Href: "https://dorm.cau.ac.kr"},
			Description: "중앙대학교 다빈치캠퍼스 기숙사의 공지사항입니다",
			Created:     time.Now(),
		}

		ctx.Set("feed", feed)
		ctx.Set("articles", articles)
	}, serveFeed)
	server.GET("/cau/dormitory/seoul/:buildingType/:feedType", func(ctx *gin.Context) {
		setRedisFeedCacheKey(ctx, ctx.Request.URL.Path)
	}, serveCachedFeed, func(ctx *gin.Context) {
		buildingType := ctx.Param("buildingType")

		var articles []cau_parser.CAUArticle
		var articlesErr error
		var buildingName string
		switch buildingType {
		case "bluemir":
			buildingName = " 블루미르홀"
		case "future_house":
			buildingName = " 퓨처하우스"
		case "global_house":
			buildingName = " 글로벌하우스"
		case "all":
			buildingName = ""
		default:
			fmt.Fprintf(os.Stderr, "unsupported building: %s", ctx.Param("buildingType"))
			ctx.HTML(404, "404.html", gin.H{})
			return
		}

		if buildingType == "all" {
			articles, articlesErr = getAllSeoulDormitoryArticles()
		} else {
			articles, articlesErr = getArticlesWithCache("dormitory/seoul/" + buildingType)
		}

		if articlesErr != nil {
			fmt.Fprint(os.Stderr, articlesErr)
			ctx.HTML(500, "article-fetch-error.html", gin.H{})
			return
		}

		feed := &feeds.Feed{
			Title:       "중앙대학교 서울캠퍼스 기숙사" + buildingName + " 공지사항",
			Link:        &feeds.Link{Href: "https://dormitory.cau.ac.kr"},
			Description: "중앙대학교 서울캠퍼스 기숙사" + buildingName + "의 공지사항입니다",
			Created:     time.Now(),
		}

		ctx.Set("feed", feed)
		ctx.Set("articles", articles)
	}, serveFeed)
	server.Static("/img", "static/img")
	server.StaticFile("/robots.txt", "static/robots.txt")

	return server
}
