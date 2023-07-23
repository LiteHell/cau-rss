package cau_parser_test

import (
	"testing"

	"litehell.info/cau-rss/cau_parser"
)

func TestParsers(t *testing.T) {
	t.Run("abeek", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseABEEK()

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("ai", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseAI()

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("cse", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseCSE()

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("dormitory (bluemir)", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_BLUEMIR)

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("dormitory (future house)", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_FUTURE_HOUSE)

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("dormitory (global house)", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_GLOBAL_HOUSE)

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("(dormitory davinci)", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseDormitory(cau_parser.DORMITORY_DAVINCI)

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("dormitory (human)", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseHuman()

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("dormitory (swedu)", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseSWEDU()

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})

	t.Run("ict", func(t *testing.T) {
		t.Parallel()
		articles, err := cau_parser.ParseICT()

		if err != nil {
			t.Error(err)
		}
		testArticles(articles, t)
	})
}
