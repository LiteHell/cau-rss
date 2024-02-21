package cau_parser

import (
	"crypto/tls"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func getHtmlFromUrl(url string) (*goquery.Document, error) {
	// Fetch, ignoring certificate error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
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
