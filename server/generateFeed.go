package server

import (
	"fmt"
	"html"
	"time"

	"github.com/gorilla/feeds"
	"litehell.info/cau-rss/cau_parser"
)

type feedType int8

const RSS = 0
const ATOM = 1
const JSON = 2

func GenerateFeed(feed *feeds.Feed, articles []cau_parser.CAUArticle, feedType feedType) (string, error) {
	feed.Image = &feeds.Image{
		Url:    getWebAddress() + "/img/puang.png",
		Title:  "RSS 마크를 껴안은 푸앙이",
		Link:   getWebAddress(),
		Width:  400,
		Height: 400,
	}

	updated := time.Unix(0, 0)

	for _, article := range articles {
		var files string = ""
		for _, file := range article.Files {
			files += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", html.EscapeString(file.Url), html.EscapeString(file.Name))
		}
		if files != "" {
			files = fmt.Sprintf("<div style=\"boder: 1px solid black; padding: 10px;\"><p>첨부파일</p><ul>%s</ul></div>", files)
		}
		feed.Add(&feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: article.Url},
			Author:      &feeds.Author{Name: article.Author},
			Created:     article.Date,
			Content:     article.Content + files,
			Description: article.Content,
			Id:          article.Url,
		})
		if updated.Before(article.Date) {
			updated = article.Date
		}
	}

	feed.Updated = updated

	var result string
	var err error
	switch feedType {
	case RSS:
		result, err = feed.ToRss()
	case ATOM:
		result, err = feed.ToAtom()
	case JSON:
		result, err = feed.ToJSON()
	default:
		return "", fmt.Errorf("Unknown feed type")
	}

	if err != nil {
		return "", err
	}
	return result, nil
}
