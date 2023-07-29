package server

import "os"

func getWebAddress() string {
	webAddress := os.Getenv("WEB_ADDRESS")
	if webAddress == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		webAddress = "http://127.0.0.1:" + port
	}

	return webAddress
}
