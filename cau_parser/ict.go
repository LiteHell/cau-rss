package cau_parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getICTArticle(articleUrl string) (string, []CAUAttachment, error) {
	// Fetch article
	html, err := getHtmlFromUrl(articleUrl)
	if err != nil {
		return "", nil, err
	}

	// Get content
	content, contentErr := html.Find("table.con-tb tbody tr > td#bo-cont").Html()
	if contentErr != nil {
		return "", nil, contentErr
	}

	// Get files
	var fileErr error
	fileHrefPattern := regexp.MustCompile("javascript:file_down\\(['\"](.+?)['\"],['\"](.+?)['\"],['\"](.+?)['\"],['\"](.+?)['\"]\\)")
	fileLinks := html.Find("table.con-tb tbody tr > td.file a")
	files := make([]CAUAttachment, fileLinks.Size())
	fileLinks.Each(func(i int, s *goquery.Selection) {
		matches := fileHrefPattern.FindStringSubmatch(s.AttrOr("href", ""))
		originalFileName, err := url.QueryUnescape(matches[4])
		if err != nil {
			fileErr = err
			return
		}

		files[i].Name = originalFileName
		files[i].Url = fmt.Sprintf("https://ict.cau.ac.kr/ModulePrint/ModuleInclude/filedown.php?m_name=%s&m_dir=%s&file_name=%s&org_file_name=%s",
			url.QueryEscape(matches[1]),
			url.QueryEscape(matches[2]),
			url.QueryEscape(matches[3]),
			url.QueryEscape(matches[4]))

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

func ParseICT() ([]CAUArticle, error) {
	boardUrl := "https://ict.cau.ac.kr/em/notice.php"
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Get rows
	rows := html.Find(".con-tb tbody tr")
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
		articles[i].Author = "창의ICT공과대학"

		content, files, contentErr := getICTArticle(articles[i].Url)
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
