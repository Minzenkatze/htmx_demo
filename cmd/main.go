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
	db, err := db.NewSqliteDB()
	defer func() {
		if db.Db != nil {
			db.Db.Close()
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	handlers := &handlers.Handlers{
		Db:     db,
		Router: r,
	}
	handlers.SetupRoutes()
	r.Run(":8080")
}
