package main

import (
	"log"

	"litehell.info/cau-rss/cau_parser"
)

func main() {
	test, err := cau_parser.ParseCSE()
	if err != nil {
		panic(err)
	}

	for _, i := range test {
		log.Println(i)
	}
}
