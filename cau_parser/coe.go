package cau_parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getCOEArticle(url string) (string, string, []CAUAttachment, error) {
	// Fetch html
	html, err := getHtmlFromUrl(url)
	if err != nil {
		return "", "", nil, err
	}

	title := html.Find(".board-view > .head > h3.tit").Text()
	content, err := html.Find(".board-view > .body").Html()
	if err != nil {
		return "", "", nil, err
	}

	fileLinks := html.Find("#fileLayer a[title]")
	files := make([]CAUAttachment, fileLinks.Size())
	// "/module/board/download.php?boardid="+boardid+"&b_idx="+b_idx+"&idx="+idx
	// javascript:download('notice','828','1420');

	pattern := regexp.MustCompile("javascript:download\\(['\"](.+?)['\"],['\"](.+?)['\"],['\"](.+?)['\"]\\)")
	fileLinks.Each(func(i int, s *goquery.Selection) {
		matches := pattern.FindStringSubmatch(s.AttrOr("href", ""))

		files[i].Name = s.AttrOr("title", s.Text())
		files[i].Url = fmt.Sprintf("https://coe.cau.ac.kr/module/board/download.php?boardid=%s&b_idx=%s&idx=%s", matches[1], matches[2], matches[3])
	})

	return title, content, files, nil
}

func ParseCOE() ([]CAUArticle, error) {
	// Fetch html
	boardUrl := "https://coe.cau.ac.kr/sub/sub04_01.php"
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Get rows
	rows := html.Find(".board-list table tbody tr")
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	var rowErr error
	rows.Each(func(i int, s *goquery.Selection) {
		subjectLink := s.Find(".subject a")
		author := s.Find(".name").Text()
		date, dateErr := time.Parse("2006-01-02", strings.TrimSpace(s.Find(".date").Text()))
		url, urlErr := getAbsoluteUrlFromRelative(strings.TrimSpace(subjectLink.AttrOr("href", "")), boardUrl)
		title, content, files, contentErr := getCOEArticle(url.String())
		if contentErr != nil || dateErr != nil || urlErr != nil {
			rowErr = contentErr
			return
		}
		articles[i].Title = strings.TrimSpace(title)
		articles[i].Url = url.String()
		articles[i].Author = strings.TrimSpace(author)
		articles[i].Date = date
		articles[i].Content = content
		articles[i].Files = files
	})

	if rowErr != nil {
		return nil, rowErr
	}

	return articles, nil
}
