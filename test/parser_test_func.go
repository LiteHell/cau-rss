package test

import (
	"strings"
	"testing"

	"litehell.info/cau-rss/cau_parser"
)

func testArticles(articles []cau_parser.CAUArticle, t *testing.T) {
	if len(articles) == 0 {
		t.Error("No articles")
	}

	for _, article := range articles {
		switch {
		case article.Author == "":
			t.Error("Empty author")
		case article.Author != strings.TrimSpace(article.Author):
			t.Error("Author not trimmed")
		case article.Title == "":
			t.Error("Empty title")
		case article.Title != strings.TrimSpace(article.Title):
			t.Error("Title not trimmed")
		case article.Content == "":
			t.Error("Empty content")
		case article.Url == "":
			t.Error("Empty url")
		}
		t.Logf("Article (Author: %s, Title: %s, Url: %s)",
			article.Author, article.Title, article.Url)
		if len(article.Files) > 0 {
			for _, file := range article.Files {
				switch {
				case file.Name == "":
					t.Error("Empty file name")
				case file.Name != strings.TrimSpace(file.Name):
					t.Error("File name not trimmed")
				case file.Url == "":
					t.Error("Empty url")
				}
				t.Logf("File (Name: %s, Url: %s)", file.Name, file.Url)
			}
		}
		//t.Logf("Content:\n\t%s", article.Content)
	}
}
