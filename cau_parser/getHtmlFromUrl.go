package cau_parser

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func getHtmlFromUrl(url string) (*goquery.Document, error) {
	// Fetch
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse html
	html, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return html, nil
}
