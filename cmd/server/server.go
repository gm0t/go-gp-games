package main

import (
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

type EvolveRequest struct {
	PopulationSize     int     `json:"populationSize"`
	ElitesSize         int     `json:"elitesSize"`
	MutationChance     float64 `json:"mutationChance"`
	NumberOfCrossovers int     `json:"numberOfCrossovers"`
}

func main() {
	generator := tree.NewGenerator(catcher.BoolKeys, catcher.FloatKeys, []string{
		string(actions.Up),
		string(actions.Down),
		string(actions.Left),
		string(actions.Right),
	})

	router := gin.Default()
	population := evolution.NewPopulation(&evolution.Params{
		Size:           50,
		ElitesSize:     5,
		ChildrenSize:   20,
		Generator:      generator,
		Fitness:        catcher.NewFitness(),
		MutationChance: 0.5,
	})

	go population.Evolve(catcher.NewTerminationCondition())

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"generation": population.CurrentGeneration(),
			"finished":   population.IsFinished(),
			"genes":      population.Genes(),
			"elites":     population.Elites(),
			"best":       population.Best(),
		})
	})

	router.POST("/evolve", func(c *gin.Context) {
		request := &EvolveRequest{}
		c.BindJSON(request)

		if request.PopulationSize < 6 {
			c.String(http.StatusBadRequest, "PopulationSize is too low")
			return
		}

		if request.ElitesSize < 0 || request.ElitesSize > request.PopulationSize {
			c.String(http.StatusBadRequest, "ElitesSize is wrong")
			return
		}

		if request.NumberOfCrossovers < 0 || request.NumberOfCrossovers > request.PopulationSize {
			c.String(http.StatusBadRequest, "NumberOfCrossovers is wrong [0, %d]", request.PopulationSize)
			return
		}

		if request.MutationChance < 0 || request.MutationChance > 1 {
			c.String(http.StatusBadRequest, "MutationChance is wrong [0, 1]")
			return
		}

		wait := make(chan interface{})
		population.Stop(wait)
		<-wait

		population = evolution.NewPopulation(&evolution.Params{
			Size:           request.PopulationSize,
			ElitesSize:     request.ElitesSize,
			ChildrenSize:   request.NumberOfCrossovers,
			Generator:      generator,
			Fitness:        catcher.NewFitness(),
			MutationChance: request.MutationChance,
		})

		go population.Evolve(catcher.NewTerminationCondition())
	})

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
			"actions":          result.Actions,
			"normalized":       normalized,
			"asString":         normalized.String(),
			"asStringOriginal": request.Agent.String(),
		})
	})
	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	router.Run(":9000")
}
