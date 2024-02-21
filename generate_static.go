package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path"

	"github.com/gorilla/feeds"
	cp "github.com/otiai10/copy"
	"litehell.info/cau-rss/server"
)

func generateIndex(dir string) {
	indexTemplate, err := template.New("").Funcs(
		template.FuncMap{
			"encodeURI": func(uri string) string {
				return url.QueryEscape(uri)
			},
		},
	).ParseFiles("html/index.html")
	if err != nil {
		panic(err)
	}

	indexOutputFile, err := os.Create(path.Join(dir, "index.html"))
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(indexOutputFile)
	indexTemplate.ExecuteTemplate(writer, "index.html", map[string]any{
		"table":      server.GetFeedHtmlTable(),
		"webAddress": "rss.puang.network",
	})

	defer func() {
		writer.Flush()
		indexOutputFile.Close()
	}()
}

func generateFeedFiles(dir string, items *FeedDataResponse) {
	for _, item := range items.Success {
		feedName := item.SiteInfo.Name
		if item.SiteInfo.LongName != "" {
			feedName = item.SiteInfo.LongName
		}
		feed := &feeds.Feed{
			Title:       fmt.Sprintf("%s 공지사항", feedName),
			Link:        &feeds.Link{Href: item.SiteInfo.Url},
			Description: fmt.Sprintf("%s의 공지사항입니다", feedName),
		}

		rss, err := server.GenerateFeed(feed, item.Articles, server.RSS)
		if err != nil {
			panic(err)
		}

		atom, err := server.GenerateFeed(feed, item.Articles, server.ATOM)
		if err != nil {
			panic(err)
		}

		json, err := server.GenerateFeed(feed, item.Articles, server.JSON)
		if err != nil {
			panic(err)
		}

		feedDir := path.Join(dir, "cau", item.SiteInfo.Key)
		err = os.MkdirAll(feedDir, 0755)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(path.Join(feedDir, "rss"), []byte(rss), 0644)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(path.Join(feedDir, "atom"), []byte(atom), 0644)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(path.Join(feedDir, "json"), []byte(json), 0644)
		if err != nil {
			panic(err)
		}

	}
}

func generateStatic(dir string, items *FeedDataResponse) {
	generateIndex(dir)
	generateFeedFiles(dir, items)
	cp.Copy("static", dir)
}
