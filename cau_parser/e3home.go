package cau_parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getE3HomeArticle(articleUrl string) (string, []CAUAttachment, error) {
	// Fetch article
	html, err := getHtmlFromUrl(articleUrl)
	if err != nil {
		return "", nil, err
	}

	// Get content
	content, contentErr := html.Find("#em_w_con1").Html()
	if contentErr != nil {
		return "", nil, contentErr
	}

	// Get files
	var fileErr error
	fileHrefPattern := regexp.MustCompile("javascript:download\\(['\"](.+?)['\"], ['\"](.+?)['\"], ['\"](.+?)['\"]\\)")
	fileLinks := html.Find(".em_w_nav1 li.btn_next a.n_file2")
	files := make([]CAUAttachment, fileLinks.Size())
	fileLinks.Each(func(i int, s *goquery.Selection) {
		matches := fileHrefPattern.FindStringSubmatch(s.AttrOr("href", ""))

		files[i].Name = matches[2]
		files[i].Url = fmt.Sprintf("https://e3home.cau.ac.kr/common/download.php?downpath=%s&org_name=%s&file_name=%s",
			url.QueryEscape(matches[1]),
			url.QueryEscape(matches[2]),
			url.QueryEscape(matches[3]),
		)

		if err != nil {
			fileErr = err
			return
		}
	})

	if fileErr != nil {
		return content, nil, fileErr
	}

	return content, files, nil
}

func ParseE3Home() ([]CAUArticle, error) {
	boardUrl := "https://e3home.cau.ac.kr/em/em_1.php"
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Get rows
	rows := html.Find(".e_tbl_wrap tbody tr")
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	var rowErr error
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		title := s.Find("a")

		date, dateErr := time.Parse("2006.01.02", strings.TrimSpace(getTextFromNode(cells.Get(2))))

		if dateErr != nil {
			rowErr = dateErr
			return
		}

		articles[i].Date = date
		articles[i].Title = strings.TrimSpace(title.Text())
		articles[i].Url = javascriptViewHrefToUrl(title.AttrOr("href", ""), boardUrl)
		articles[i].Author = "전기전자공학부"

		content, files, contentErr := getE3HomeArticle(articles[i].Url)
		if contentErr != nil {
			return
		}
		articles[i].Files = files
		articles[i].Content = content
	})
	if rowErr != nil {
		return nil, rowErr
	}

	return articles, nil
}
