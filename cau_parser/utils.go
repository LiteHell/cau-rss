package cau_parser

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func getAbsoluteUrlFromRelative(relative string, base string) (*url.URL, error) {
	baseParsed, baseErr := url.Parse(base)
	if baseErr != nil {
		return nil, baseErr
	}

	relativeParsed, relativeErr := url.Parse(relative)
	if relativeErr != nil {
		return nil, relativeErr
	}

	return baseParsed.ResolveReference(relativeParsed), nil
}

func getTextFromNode(node *html.Node) string {
	return goquery.NewDocumentFromNode(node).Text()
}
