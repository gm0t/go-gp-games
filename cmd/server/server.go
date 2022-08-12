package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"lr1Go/pkg/evolution"
	"lr1Go/pkg/game/catcher"
	"lr1Go/pkg/game/catcher/actions"
	"lr1Go/pkg/tree"
)

type SimulateRequest struct {
	Steps   int       `json:"steps"`
	Agent   tree.Node `json:"agent"`
	AgentX  float64   `json:"agentX"`
	AgentY  float64   `json:"agentY"`
	TargetX float64   `json:"targetX"`
	TargetY float64   `json:"targetY"`
}

func main() {
	generator := tree.NewGenerator(catcher.BoolKeys, catcher.FloatKeys, []string{
		string(actions.Up),
		string(actions.Down),
		string(actions.Left),
		string(actions.Right),
	})

	router := gin.Default()
	population := evolution.NewPopulation(50, generator, catcher.NewFitness(), 0.5)

	go func() {
		population.Evolve(150000)
	}()

	router.GET("/status", func(c *gin.Context) {
		fmt.Println(gin.H{
			"generation": population.CurrentGeneration(),
			"best":       population.Best(),
		})

		c.JSON(http.StatusOK, gin.H{
			"generation": population.CurrentGeneration(),
			"genes":      population.Genes(),
			"best":       population.Best(),
		})
	})

	// For each matched request Context will hold the route definition
	router.POST("/catcher", func(c *gin.Context) {
		request := &SimulateRequest{}
		c.BindJSON(request)

		game := catcher.NewGame(
			catcher.NewState(
				request.AgentX,
				request.AgentY,
				request.TargetX,
				request.TargetY,
			),
			catcher.NewAiFPlayer(&request.Agent),
		)

		normalized := catcher.Normalize(&request.Agent)

		result := game.Run(request.Steps)

		c.JSON(http.StatusOK, gin.H{
			"actions":    result.Actions,
			"normalized": normalized,
			"asString":   normalized.String(),
		})
	})
	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	router.Run(":9000")
}
