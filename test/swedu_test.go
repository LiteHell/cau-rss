package test

import (
	"testing"

	"litehell.info/cau-rss/cau_parser"
)

func TestSWEDU(t *testing.T) {
	articles, err := cau_parser.ParseSWEDU()

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}
