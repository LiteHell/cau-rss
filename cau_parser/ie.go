package cau_parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getIEArticle(articleUrl string) (string, []CAUAttachment, error) {
	// Fetch article
	html, err := getHtmlFromUrl(articleUrl)
	if err != nil {
		return "", nil, err
	}

	// Get file and content container
	containers := html.Find("tr td.ali-left")

	// Get files
	files := []CAUAttachment{}
	if containers.Size() >= 2 {
		var fileErr error
		fileHrefPattern := regexp.MustCompile("javascript:file_down\\(['\"](.+?)['\"],['\"](.+?)['\"],['\"](.+?)['\"],['\"](.+?)['\"]\\)")
		fileContainer := goquery.NewDocumentFromNode(containers.Get(0))
		fileLinks := fileContainer.Find("a")
		files = make([]CAUAttachment, fileLinks.Size())
		fileLinks.Each(func(i int, s *goquery.Selection) {
			matches := fileHrefPattern.FindStringSubmatch(s.AttrOr("href", ""))

			files[i].Name = matches[2]
			files[i].Url = fmt.Sprintf("https://ie.cau.ac.kr/common/download.php?m_name=%s&m_dir=%s&file_name=%s&org_file_name=%s",
				url.QueryEscape(matches[1]),
				url.QueryEscape(matches[2]),
				url.QueryEscape(matches[3]),
				url.QueryEscape(matches[4]),
			)

			if err != nil {
				fileErr = err
				return
			}
		})

		if fileErr != nil {
			return "", nil, fileErr
		}

		fileContainer.Remove()
	}

	// Get content
	content, contentErr := goquery.NewDocumentFromNode(containers.Get(0)).Html()
	if contentErr != nil {
		return "", nil, contentErr
	}

	return content, files, nil
}

func ParseIE() ([]CAUArticle, error) {
	boardUrl := "https://ie.cau.ac.kr/sch_5/notice.php"
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Get rows
	rows := html.Find(".table-notice tbody tr")
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	var rowErr error
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		title := s.Find("a")

		date, dateErr := time.Parse("2006-01-02", strings.TrimSpace(getTextFromNode(cells.Get(4))))

		if dateErr != nil {
			date, dateErr = time.Parse("2006.01.02", strings.TrimSpace(getTextFromNode(cells.Get(4))))
			if dateErr != nil {
				rowErr = dateErr
				return
			}
		}

		articles[i].Date = date
		articles[i].Title = strings.TrimSpace(title.Text())
		articles[i].Url = javascriptViewHrefToUrl(title.AttrOr("href", ""), boardUrl)
		articles[i].Author = "융합공학부"

		content, files, contentErr := getIEArticle(articles[i].Url)
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
