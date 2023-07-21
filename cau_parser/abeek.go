package cau_parser

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseABEEKArticle(url string) (string, string, []CAUAttachment, error) {
	// Fetch article
	resp, err := http.Get(url)
	if err != nil {
		return "", "", nil, err
	}

	// Parse html
	html, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", "", nil, err
	}

	// Get content
	content, contentErr := html.Find("#contents table tr.view_text table td").Html()
	if contentErr != nil {
		return "", "", nil, contentErr
	}

	// Get nodes for author and files
	viewListCells := html.Find("#contents table tr.view_list")
	files := make([]CAUAttachment, viewListCells.Size()-1)
	var author string
	var cellErr error

	// Process author and files
	viewListCells.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// Parse author
			author = strings.TrimSpace(getTextFromNode(s.Find("td").Get(0)))
		} else {
			i--
			// Check whether it's attahcment
			if strings.TrimSpace(s.Find("th").Text()) != "첨부파일" {
				return
			}

			// Process attachment
			fileLink := s.Find("td a")
			url, urlErr := getAbsoluteUrlFromRelative(fileLink.AttrOr("href", ""), url)
			if urlErr != nil {
				cellErr = err
				return
			}

			files[i].Name = strings.TrimSpace(fileLink.Text())
			files[i].Url = url.String()
		}
	})

	if cellErr != nil {
		return "", "", nil, cellErr
	}

	// Process empty files
	var nonEmptyFiles []CAUAttachment = []CAUAttachment{}
	for _, file := range files {
		if file.Name != "" && file.Url != "" {
			nonEmptyFiles = append(nonEmptyFiles, file)
		}
	}

	return content, author, nonEmptyFiles, nil
}

func ParseABEEK() ([]CAUArticle, error) {
	// Fetch board
	boardUrl := "https://abeek.cau.ac.kr/notice/list.jsp?sc_board_seq=1&page=1"
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
	rows := html.Find("#contents table tbody tr")
	articles := make([]CAUArticle, rows.Size())

	// Variable to handle error in processing rows
	var rowErr error

	// Process rows
	rows.Each(func(i int, s *goquery.Selection) {
		date, dateErr := time.Parse(
			"2006.01.02", strings.TrimSpace(getTextFromNode(s.Find("td").Get(3))))
		if dateErr != nil {
			rowErr = dateErr
			return
		}

		articles[i].Url = fmt.Sprintf("https://abeek.cau.ac.kr/notice/view.jsp?sc_board_seq=1&pk_seq=%s",
			s.Find("td.left a").AttrOr("seq", ""))
		articles[i].Title = strings.TrimSpace(s.Find("td.left a").Text())
		articles[i].Date = date.Add(time.Hour * 9) // KST

		content, author, files, articleParseErr := parseABEEKArticle(articles[i].Url)
		if articleParseErr != nil {
			rowErr = articleParseErr
			return
		}

		articles[i].Content = content
		articles[i].Author = author
		articles[i].Files = files
	})

	// Return result
	if rowErr != nil {
		return nil, rowErr
	}
	return articles, nil
}
