package test

import (
	"testing"

	"litehell.info/cau-rss/cau_parser"
)

func TestABEEK(t *testing.T) {
	articles, err := cau_parser.ParseABEEK()

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}
