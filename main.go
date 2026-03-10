package main

import (
	"fmt"
	"mygo/config"
	"mygo/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectinDB()
	config.InitRedis()

	r := gin.Default()
	routes.SetupRoutes(r)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Golang Web Server is running!")
	})

	r.Run(":6000")
}
