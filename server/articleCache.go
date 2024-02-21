package server

import (
	"fmt"

	"litehell.info/cau-rss/cau_parser"
)

func FetchArticlesForKey(key string) ([]cau_parser.CAUArticle, error) {
	switch key {
	case "cse":
		return cau_parser.ParseCSE()
	case "swedu":
		return cau_parser.ParseSWEDU()
	case "abeek":
		return cau_parser.ParseABEEK()
	case "dormitory/davinci":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_DAVINCI)
	case "dormitory/seoul/bluemir":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_BLUEMIR)
	case "dormitory/seoul/future_house":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_FUTURE_HOUSE)
	case "dormitory/seoul/global_house":
		return cau_parser.ParseDormitory(cau_parser.DORMITORY_GLOBAL_HOUSE)
	case "ie":
		return cau_parser.ParseIE()
	case "e3home":
		return cau_parser.ParseE3Home()
	case "ict":
		return cau_parser.ParseICT()
	case "ai":
		return cau_parser.ParseAI()
	default:
		panic(fmt.Errorf("Unknown articles key: %s", key))
	}
}
