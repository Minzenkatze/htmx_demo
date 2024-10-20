package main

import (
	"htmx_demo/internal/db"
	"htmx_demo/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Static("/static", "./static")
	r.StaticFile("/styles.css", "./internal/assets/dist/styles.css")
	err := db.InitDB()
	defer func() {
		if db.DB != nil {
			db.DB.Close()
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	handlers.SetupRoutes(r)
	r.Run(":8080")
}
