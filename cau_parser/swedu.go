package cau_parser

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ParseSWEDUArticle(url string) (string, []CAUAttachment, error) {
	// Fetch article
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}

	// Parse html
	html, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil, err
	}

	// Get file links
	fileLinks := html.Find("#boardForm table tbody tr:first-child td[colspan] a")
	files := make([]CAUAttachment, fileLinks.Size())

	// Process file links
	var linkErr error
	fileLinks.Each(func(i int, s *goquery.Selection) {
		url, urlErr := getAbsoluteUrlFromRelative(s.AttrOr("href", ""), url)
		if urlErr != nil {
			linkErr = urlErr
			return
		}

		files[i].Name = strings.TrimSpace(s.Text())
		files[i].Url = url.String()
	})

	// Get content
	content, contentErr := html.Find("#boardForm table tbody td.editTd").Html()

	// Handle errors
	if linkErr != nil {
		return "", nil, linkErr
	} else if contentErr != nil {
		return "", nil, contentErr
	}

	// Return result
	return content, files, nil
}

func ParseSWEDU() ([]CAUArticle, error) {
	// Fetch board
	const boardUrl = "https://swedu.cau.ac.kr/board/list?boardtypeid=7&menuid=001005005"
	resp, err := http.Get(boardUrl)
	if err != nil {
		return nil, err
	}

	// Parse html
	html, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	// Get article rows
	rows := html.Find("#boardForm table tbody tr")
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	var rowErr error
	rows.Each(func(i int, s *goquery.Selection) {
		url, urlErr := getAbsoluteUrlFromRelative(s.Find("td.tl a").AttrOr("href", ""), "https://swedu.cau.ac.kr/board/list")
		date, dateErr := time.Parse("2006-01-02",
			strings.TrimSpace(getTextFromNode(s.Find("td").Get(2))))

		if urlErr != nil {
			rowErr = urlErr
			return
		} else if dateErr != nil {
			rowErr = dateErr
			return
		}

		articles[i].Title = strings.TrimSpace(s.Find("td.tl a").Text())
		articles[i].Url = url.String()
		articles[i].Author = "다빈치SW교육원"
		articles[i].Date = date.Add(time.Hour * 9) // KST

		content, files, articleErr := ParseSWEDUArticle(articles[i].Url)
		if articleErr != nil {
			rowErr = articleErr
			return
		}

		articles[i].Content = content
		articles[i].Files = files
	})

	if rowErr != nil {
		return nil, rowErr
	}
	return articles, nil
}
