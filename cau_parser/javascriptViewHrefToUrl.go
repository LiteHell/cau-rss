package cau_parser

import "regexp"

var javascriptViewPattern = regexp.MustCompile("javascript:view\\([\"']?([0-9]+)")

func javascriptViewHrefToUrl(jsHref, baseUrl string) string {
	id := javascriptViewPattern.FindStringSubmatch(jsHref)[1]
	return baseUrl + "?p_mode=view&p_idx=" + id
}
