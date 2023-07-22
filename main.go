package main

import (
	"litehell.info/cau-rss/server"
)

func main() {
	server.InitializeRedis()

	server := server.CreateServer()
	server.Run()
}
