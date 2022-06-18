package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lr1Go/pkg/old-tree"
)

type SimulateRequest struct {
	MaxAttempts int                 `json:"maxAttempts"`
	Agent       old_tree.Serialized `json:"agent"`
}

func main() {

	router := gin.Default()

	router.GET("/evolve", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"generation": 0, "population": []int{}})
	})
	router.GET("/simulate", func(c *gin.Context) {
		req := &SimulateRequest{}
		c.BindJSON(req)

		c.JSON(http.StatusOK, gin.H{"generation": 0, "population": []int{}})
	})

	// For each matched request Context will hold the route definition
	router.POST("/evolve", func(c *gin.Context) {
		b := c.FullPath() == "/user/:name/*action" // true
		c.String(http.StatusOK, "%t", b)
	})
	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	router.Run(":9000")
}
