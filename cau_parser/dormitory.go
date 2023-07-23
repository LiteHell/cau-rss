package cau_parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const DORMITORY_BLUEMIR = "m05_01_01"
const DORMITORY_FUTURE_HOUSE = "m05_01_02"
const DORMITORY_GLOBAL_HOUSE = "m05_01_03"
const DORMITORY_DAVINCI = "m06_01"

func parseDormitoryArticle(url string) (string, []CAUAttachment, error) {
	// Fetch board
	html, err := getHtmlFromUrl(url)
	if err != nil {
		return "", nil, err
	}

	// There are files?
	hasFiles := html.Find("#BoardViewAdd").Size() > 1
	files := []CAUAttachment{}

	// Process files if exists
	if hasFiles {
		// Get file link container
		fileContainer := goquery.NewDocumentFromNode(html.Find("#BoardViewAdd").Get(1))

		// Delete images
		fileContainer.Find("img").Remove()

		// Process file links
		var fileUrlErr error
		fileLinks := fileContainer.Find("a")
		files = make([]CAUAttachment, fileLinks.Size())
		fileLinks.Each(func(i int, s *goquery.Selection) {
			url, urlErr := getAbsoluteUrlFromRelative(s.AttrOr("href", ""), url)

			if urlErr != nil {
				fileUrlErr = urlErr
				return
			}

			files[i].Name = strings.TrimSpace(s.Text())
			files[i].Url = url.String()
		})

		if fileUrlErr != nil {
			return "", nil, fileUrlErr
		}
	}

	// Get content
	content, contentErr := html.Find("#BoardContent").Html()
	if contentErr != nil {
		return "", nil, contentErr
	}

	return content, files, nil
}

func ParseDormitory(boardId string) ([]CAUArticle, error) {
	// Fetch board
	boardUrl := "https://dormitory.cau.ac.kr/community.php?mid=" + boardId
	if boardId == DORMITORY_DAVINCI {
		boardUrl = "https://dorm.cau.ac.kr/community.php?mid=m06_01"
	}
	html, err := getHtmlFromUrl(boardUrl)
	if err != nil {
		return nil, err
	}

	// Delete NEW, HOT tags
	html.Find("table#Board img").Remove()

	// Get rows
	rows := html.Find("table#Board tbody tr")
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	var rowErr error
	rows.Each(func(i int, s *goquery.Selection) {
		titleLinkNode := goquery.NewDocumentFromNode(s.Find("td.Subject > a").Get(0))
		date, dateErr := time.Parse(
			"2006-01-02",
			strings.TrimSpace(s.Find("td.board_date").Text()))
		url, urlErr := getAbsoluteUrlFromRelative(titleLinkNode.AttrOr("href", ""), boardUrl)

		if dateErr != nil {
			rowErr = dateErr
			return
		} else if urlErr != nil {
			rowErr = urlErr
			return
		}

		category := s.Find("td.Subject > span").Text()
		articles[i].Date = date
		articles[i].Title = strings.TrimSpace(fmt.Sprintf("%s %s", category, strings.TrimSpace(titleLinkNode.Text())))
		articles[i].Url = url.String()
		articles[i].Author = "기숙사"

		content, files, articleErr := parseDormitoryArticle(url.String())
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
