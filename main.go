package main

import (
	"log"
	"strings"
)

func main() {
	c := getConfig()

	app, err := build(c)
	if err != nil {
		log.Fatal("Error building app:", err)
	}

	log.Println("Loaded database", c.database)

	var url string
	if strings.HasPrefix(app.config.address, ":") {
		url = "http://localhost" + app.config.address
	} else {
		url = "http://" + app.config.address
	}
	url = url + c.root

	log.Println("Launching server on", url)

	app.run()
}
