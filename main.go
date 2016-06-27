package main

import (
	"flag"
	"restdb/router"
)

var (
	dataBaseFile string
	dataBaseName string
	address string
)

func init() {
	flag.StringVar(&dataBaseFile, "f", "./db/base.db", "DataBase file")
	flag.StringVar(&dataBaseName, "n", "users", "DataBase name")
	flag.StringVar(&address, "a", ":8080", "Network Address")
}

func main() {
	flag.Parse()

	router := router.GetRouter(dataBaseFile)

	router.Run(address)
}