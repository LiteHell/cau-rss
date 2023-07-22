package main

import (
	"litehell.info/cau-rss/server"
)

func main() {
	server.InitializeRedis()
	server.StartCrawller()

	server := server.CreateServer()
	server.Run()
}
