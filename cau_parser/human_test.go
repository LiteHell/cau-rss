package cau_parser_test

import (
	"testing"

	"litehell.info/cau-rss/cau_parser"
)

func TestHuman(t *testing.T) {
	articles, err := cau_parser.ParseHuman()

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}
