package main

import (
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRouter(endpointMap map[string]Endpoint, ready *bool) (r *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if !*ready {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"healthy": false,
				"cause":   "not ready yet",
			})
			return
		}

		// Add other checks here as necessary
		c.JSON(http.StatusOK, gin.H{
			"healthy": true,
		})
	})
	r.GET("/ready", func(c *gin.Context) {
		status := http.StatusServiceUnavailable
		if *ready {
			status = http.StatusOK
		}
		c.JSON(status, gin.H{
			"ready": ready,
		})
	})
	r.POST("/webhook", func(context *gin.Context) {
		body, err := context.GetRawData();if err!=nil{
			context.JSON(http.StatusBadRequest, "Error")
		}
		jsonParsed, err := gabs.ParseJSON(body)
		intent, _ := jsonParsed.Path("queryResult.intent.displayName").Data().(string)
		if intent == "spot.distance" {
			jsonResponse := gabs.New()
			jsonResponse.Set("calculating Location", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		} else if intent == "spot.available" {
			jsonResponse := gabs.New()
			jsonResponse.Set("Getting Availability", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		}else if intent == "spot.setstore" {
			jsonResponse := gabs.New()
			jsonResponse.Set("Nearest Store Set", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		}
	})

	return r
}
