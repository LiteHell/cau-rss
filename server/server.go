package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func CreateServer() *gin.Engine {
	server := gin.Default()
	server.LoadHTMLGlob("html/*.html")

	server.GET("/index.html", func(ctx *gin.Context) {
		ctx.Redirect(308, "/")
	})
	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"table": getFeedHtmlTable(),
		})
	})
	server.GET("/cau/notice", func(ctx *gin.Context) {
		ctx.Redirect(308, "https://www.cau.ac.kr/cms/FR_PRO_CON/BoardRss.do?pageNo=1&pagePerCnt=15&MENU_ID=100&SITE_NO=2&BOARD_SEQ=4&S_CATE_SEQ=&BOARD_TYPE=C0301&BOARD_CATEGORY_NO=&P_TAB_NO=&TAB_NO=&P_CATE_SEQ=&CATE_SEQ=&SEARCH_FLD=SUBJECT&SEARCH=")
	})
	for from, to := range map[string]string{"sw": "/cau/swedu/", "dormitory": "/cau/dormitory/seoul/all/"} {
		from, to := from, to
		server.GET(fmt.Sprintf("/cau/%s/:feedType", from), func(ctx *gin.Context) {
			ctx.Redirect(308, to+ctx.Param("feedType"))
		})
	}
	LoopForAllSites(func(cw *CauWebsite) {
		server.GET("/cau/"+cw.Key+"/:feedType", func(ctx *gin.Context) {
			setRedisFeedCacheKey(ctx, ctx.Request.URL.Path)
		}, serveCachedFeed, func(ctx *gin.Context) {
			feedName := cw.Name
			if cw.LongName != "" {
				feedName = cw.LongName
			}
			feed := &feeds.Feed{
				Title:       fmt.Sprintf("%s 공지사항", feedName),
				Link:        &feeds.Link{Href: cw.Url},
				Description: fmt.Sprintf("%s의 공지사항입니다", feedName),
			}

			articles, articlesErr := getArticlesWithCache(cw.Key)
			if articlesErr != nil {
				fmt.Fprint(os.Stderr, articlesErr)
				ctx.HTML(500, "article-fetch-error.html", gin.H{})
				return
			}

			ctx.Set("feed", feed)
			ctx.Set("articles", articles)
		}, serveFeed)
	})

	server.Static("/img", "static/img")
	server.StaticFile("/robots.txt", "static/robots.txt")
	server.StaticFile("/favicon.ico", "static/favicon.ico")

	return server
}
