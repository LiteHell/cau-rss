package cau_parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseABEEKArticle(url string) (string, string, []CAUAttachment, error) {
	html, err := getHtmlFromUrl(url)
	if err != nil {
		return "", "", nil, err
	}

	// Get content
	content, contentErr := html.Find("table.tb1 tbody tr:first-child td#bo-cont").Html()
	if contentErr != nil {
		return "", "", nil, contentErr
	}

	// Get nodes for author and files
	fileLinks := html.Find("td.file a")
	files := make([]CAUAttachment, fileLinks.Size())
	var cellErr error

	// Process files
	fileLinks.Each(func(i int, s *goquery.Selection) {
		resolve, resolveErr := getAbsoluteUrlFromRelative(s.AttrOr("href", ""), url)

		if resolveErr != nil {
			cellErr = resolveErr
		} else {
			files[i] = CAUAttachment{
				Name: strings.TrimSpace(s.Text()),
				Url:  resolve.String(),
			}
		}
	})

	if cellErr != nil {
		return "", "", nil, cellErr
	}

	return content, "ABEEK", files, nil
}

func ParseABEEK() ([]CAUArticle, error) {
	// Fetch board
	boardUrl := "https://abeek.cau.ac.kr/em/notice.jsp"
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Get article rows
	rows := html.Find("table.tb2 tbody tr")
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

		articles[i].Url = fmt.Sprintf("https://abeek.cau.ac.kr/em/view.jsp?sc_board_seq=1&pk_seq=%s",
			s.Find("a.btnView[seq]").AttrOr("seq", ""))
		articles[i].Title = strings.TrimSpace(s.Find("a.btnView[seq]").Text())
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
