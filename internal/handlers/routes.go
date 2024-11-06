package handlers

import (
	"htmx_demo/internal/components"
	"htmx_demo/internal/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Db     db.DB
	Router *gin.Engine
}

func (h *Handlers) renderIndex(c *gin.Context) {
	result, err := h.Db.QueryOptions()
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error querying the database")
		return
	}
	components.IndexTemplate(result).Render(c, c.Writer)
}

func (h *Handlers) filter(c *gin.Context) {
	filter := components.Filter{
		Species: c.Query("species"),
		Type:    c.Query("type"),
		Name:    c.Query("name"),
	}
	results, err := h.Db.QueryFiltered(filter)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error querying the database")
		return
	}
	components.GroupTemplate(results).Render(c, c.Writer)
}

func (h *Handlers) SetupRoutes() {
	h.Router.GET("/", h.renderIndex)
	h.Router.GET("/filter", h.filter)
}
