package cau_parser

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getCSEArticle(url string) (string, []CAUAttachment, error) {
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

	// Count files
	fileNodes := html.Find("section#content .detail .files span")
	files := make([]CAUAttachment, fileNodes.Size())

	// Parse files
	var nodeErr error = nil
	fileNodes.Each(func(i int, node *goquery.Selection) {
		// Transform file url
		url, err := getAbsoluteUrlFromRelative(regexp.MustCompile("goLocation\\(['\"](.+?)['\"] ?, ?['\"](.+?)['\"] ?, ?['\"](.+?)['\"]\\)").ReplaceAllString(
			node.AttrOr("onclick", ""), "$1?uid=$2&code=$3",
		), url)

		if err != nil {
			nodeErr = err
		}

		files[i].Url = url.String()
		files[i].Name = strings.TrimSpace(node.Text())

	})

	// Return error in parsing if exists
	if nodeErr != nil {
		return "", nil, nodeErr
	}

	// Delete attachment container from html
	html.Find("section#content .detail .files").Remove()

	// Get article content
	content, err := html.Find("section#content .detail").Html()

	// Return result
	if err != nil {
		return "", nil, nodeErr
	}
	return content, files, nil
}

func ParseCSE() ([]CAUArticle, error) {
	// Fetch cse board
	boardUrl := "http://cse.cau.ac.kr/sub05/sub0501.php"
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

	// Delete NEW tags
	html.Find(".tag.new").Remove()

	// Get rows
	var rowErr error = nil
	rows := html.Find("#listpage_form tbody tr")

	// Create article variable
	articles := make([]CAUArticle, rows.Size())

	// Process rows
	rows.Each(func(i int, row *goquery.Selection) {
		// Get td.pc-only cells for title and author
		pcOnlyCells := row.Find("td.pc-only")
		// Get article url
		articleUrl, err := getAbsoluteUrlFromRelative(
			row.Find("td.aleft a").AttrOr("href", ""),
			boardUrl)
		// Get title, author, date
		title := row.Find("td.aleft a").Text()
		author := getTextFromNode(pcOnlyCells.Get(1))
		date, dateErr := time.Parse("2006.01.02",
			getTextFromNode(pcOnlyCells.Get(2)))

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
		content, files, err := getCSEArticle(articleUrl.String())
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
