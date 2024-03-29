package main

import (
	"os"
	"./service"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main(){
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "port number for httpd service")
	if len(*pPort) != 0 {
		port = *pPort
	}

	server := service.NewServer()
	server.Run(":" + port)
}