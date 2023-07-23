package cau_parser

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getHumanArticle(url string) (string, []CAUAttachment, error) {
	html, err := getHtmlFromUrl(url)
	if err != nil {
		return "", nil, err
	}

	// 3rd row=content, 4th row=attachments
	rows := html.Find("#tbl_list_new_ct tbody tr")

	// Get content
	content, contentErr := goquery.NewDocumentFromNode(rows.Get(2)).Find("td > .bd_content").Html()
	if contentErr != nil {
		return "", nil, contentErr
	}

	// attachments
	attachmentHrefPattern := regexp.MustCompile("javascript:FileDown\\((.+?)\\)")
	attachmentLinks := goquery.NewDocumentFromNode(rows.Get(3)).Find("a")
	files := make([]CAUAttachment, attachmentLinks.Size())
	attachmentLinks.Each(func(i int, s *goquery.Selection) {
		matches := attachmentHrefPattern.FindStringSubmatch(s.AttrOr("href", ""))
		files[i].Name = strings.TrimSpace(s.Text())
		files[i].Url = "https://human.cau.ac.kr/include/FileDown.asp?FileDown=" + matches[1]
	})

	return content, files, nil
}

func ParseHuman() ([]CAUArticle, error) {
	// Fetch board
	boardUrl := "https://human.cau.ac.kr/community/notice/List.asp"
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Get rows
	rows := html.Find("#tbl_list_new_ct tbody tr")
	articles := make([]CAUArticle, rows.Size()-1)

	// Process rows
	var rowErr error
	rows.Each(func(i int, s *goquery.Selection) {
		// Skip header
		if i == 0 {
			return
		}
		i--

		// Get cells
		cells := s.Find("td")

		// Get title
		titleLink := s.Find("td > a")

		// Parse date
		date, dateErr := time.Parse("2006.01.02", strings.TrimSpace(
			getTextFromNode(cells.Get(3)),
		))
		if dateErr != nil {
			rowErr = dateErr
			return
		}

		// parse Url
		url, urlErr := getAbsoluteUrlFromRelative(titleLink.AttrOr("href", ""), boardUrl)
		if urlErr != nil {
			rowErr = urlErr
			return
		}

		articles[i].Author = strings.TrimSpace(getTextFromNode(cells.Get(2)))
		articles[i].Date = date
		articles[i].Title = titleLink.Text()
		articles[i].Url = url.String()

		content, files, contentErr := getHumanArticle(url.String())

		if contentErr != nil {
			rowErr = contentErr
			return
		}

		articles[i].Content = content
		articles[i].Files = files
	})

	if rowErr != nil {
		return nil, rowErr
	}

	return articles, rowErr
}
