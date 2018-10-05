package main

import (
	"log"
)

func main() {
	c := getConfig()

	app, err := build(c)
	if err != nil {
		log.Fatal("error building app: ", err)
	}

	app.run()
}
