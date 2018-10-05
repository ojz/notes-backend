package main

import "flag"

// config holds settings that are set by the sysadmin running your application.
type config struct {
	address  string
	database string
	dev      bool
	init     bool
}

func getConfig() config {
	c := config{}
	flag.StringVar(&c.address, "address", "localhost:8085", "The address that the server will listen on.")
	flag.BoolVar(&c.init, "init", false, "Initialize a new database.")
	flag.StringVar(&c.database, "database", "", "Path to an sqlite3 database file.")
	flag.BoolVar(&c.dev, "dev", false, "Turn on development mode, do not use this in production.")
	flag.Parse()
	return c
}
