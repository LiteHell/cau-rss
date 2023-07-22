package cau_parser

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getAIArticle(url string) (string, []CAUAttachment, error) {
	// Fetch HTTP response from cse board
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// Parse html
	html, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", nil, err
	}

	// Get content
	contentHtml, err := html.Find(".width-center.board .fr-view.detail").Html()
	if err != nil {
		return "", nil, err
	}

	// Get attachments
	fileLinks := html.Find(".width-center.board .fr-view.detail .fr-file")
	files := make([]CAUAttachment, fileLinks.Size())

	fileLinks.Each(func(i int, s *goquery.Selection) {
		files[i].Name = strings.TrimSpace(s.Text())
		files[i].Url = s.AttrOr("href", "")
	})

	return contentHtml, files, nil
}

func ParseAI() ([]CAUArticle, error) {
	// Fetch cse board
	boardUrl := "https://ai.cau.ac.kr/sub07/sub0701.php?category=1&view=list"
	resp, err := http.Get(boardUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse html
	html, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	// Get rows
	var rowErr error = nil
	rows := html.Find(".table-basic tbody tr")

	// Create article variable
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	rows.Each(func(i int, row *goquery.Selection) {
		// Get td.pc-only cells for title and author
		pcOnlyCells := row.Find("td.pc-only")
		// Get article url
		articleUrl, err := getAbsoluteUrlFromRelative(
			row.Find("td.title a").AttrOr("href", ""),
			boardUrl)

		// Get title, author, date
		title := strings.TrimSpace(row.Find("td.title a").Text())
		author := strings.TrimSpace(getTextFromNode(pcOnlyCells.Get(0)))
		date, dateErr := time.Parse("2006-01-02 15:04:05",
			strings.TrimSpace(getTextFromNode(pcOnlyCells.Get(1))))

		// Handle error
		if dateErr != nil {
			rowErr = dateErr
		} else if err != nil {
			rowErr = err
		}

		// Set article property
		articles[i].Url = articleUrl.String()
		articles[i].Title = strings.TrimSpace(title)
		articles[i].Author = strings.TrimSpace(author)
		articles[i].Date = date.Add(time.Hour * 9) // KST

		// Set content and attachments
		content, files, err := getAIArticle(articleUrl.String())
		if err != nil {
			rowErr = err
		} else {
			articles[i].Content = content
			articles[i].Files = files
		}
	})

	if rowErr != nil {
		return nil, rowErr
	}

	return articles, nil
}
