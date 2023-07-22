package server_test

import (
	"os"
	"strings"
	"time"

	"litehell.info/cau-rss/server"
)

var running_server bool = false

func runServer() {
	if running_server {
		return
	} else {
		running_server = true
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(wd, "/server") {
		os.Chdir(wd[:len(wd)-7])
	}
	server := server.CreateServer()

	go server.Run()
	time.Sleep(time.Second * 1)
}
