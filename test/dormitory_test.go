package test

import (
	"testing"

	"litehell.info/cau-rss/cau_parser"
)

func TestDormitoryBluemir(t *testing.T) {
	articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_BLUEMIR)

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}

func TestDormitoryFutureHouse(t *testing.T) {
	articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_FUTURE_HOUSE)

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}

func TestDormitoryGlobalHouse(t *testing.T) {
	articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_GLOBAL_HOUSE)

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}

func TestDormitoryDavinci(t *testing.T) {
	articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_DAVINCI)

	if err != nil {
		t.Error(err)
	}
	testArticles(articles, t)
}
