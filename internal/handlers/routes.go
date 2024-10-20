package handlers

import (
	"htmx_demo/internal/components"
	"htmx_demo/internal/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderIndex(c *gin.Context) {
	result, err := db.QueryOptions()
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error querying the database")
		return
	}
	components.IndexTemplate(result).Render(c, c.Writer)
}

func filter(c *gin.Context) {
	filter := components.Filter{
		Species: c.Query("species"),
		Type:    c.Query("type"),
		Name:    c.Query("name"),
	}
	results, err := db.QueryFiltered(filter)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error querying the database")
		return
	}
	components.GroupTemplate(results).Render(c, c.Writer)
}

func SetupRoutes(r *gin.Engine) {
	r.GET("/", renderIndex)
	r.GET("/filter", filter)
}
