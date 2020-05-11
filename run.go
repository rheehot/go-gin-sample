package main

import (
	"github.com/yoonje/gin-sample-server/server"
	"log"
)
func main() {
	log.Fatal(server.GenerateApp().Run(":5000"))
}
