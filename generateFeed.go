package main

import (
	"fmt"

	"github.com/gorilla/feeds"
	"litehell.info/cau-rss/cau_parser"
)

type feedType int8

const RSS = 0
const ATOM = 1
const JSON = 2

func generateFeed(feed *feeds.Feed, articles []cau_parser.CAUArticle, feedType feedType) (string, error) {
	for _, article := range articles {
		feed.Add(&feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: article.Url},
			Author:      &feeds.Author{Name: article.Author},
			Created:     article.Date,
			Content:     article.Content,
			Description: article.Content,
		})
	}

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
