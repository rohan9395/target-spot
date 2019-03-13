package controller

import (
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/target-spot/config"
	"github.com/target-spot/util"
	"net/http"
	"github.com/target-spot/colorlizard"
	"github.com/target-spot/store-details"
	)


func GetRouter(endpointMap map[string]config.Endpoint, ready *bool) (r *gin.Engine) {
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
		body, err := context.GetRawData()
		if err != nil {
			context.JSON(http.StatusBadRequest, "Error")
		}
		jsonParsed, err := gabs.ParseJSON(body)
		intent, _ := jsonParsed.Path("queryResult.intent.displayName").Data().(string)

		//There For Future Reference
		//contextName,contextMap := util2.ContextGet(*jsonParsed)
		//fmt.Println(contextName)
		//fmt.Print(contextMap)
		//json:=gabs.New()
		//jsonContext := util2.ContextSet(*json,"90",contextName,contextMap)
		//fmt.Print(jsonContext.String())

		switch intent {
		case "spot.distance":
			jsonResponse := gabs.New()
			jsonResponse.Set("calculating Location", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.available":
			jsonResponse := gabs.New()
			jsonResponse.Set("Getting Availability", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.setstore":
			contextName,contextMap := util.ContextGet(*jsonParsed)

			store := store_details.GetStoreID(contextMap["geo-city"].Data().(string))
			storename := store_details.GetStoreName(store)
			jsonResponse := gabs.New()
			temp := gabs.New()
			temp.Set(store)
			temp1 := gabs.New()
			temp1.Set(storename)

			contextMap["store"] = temp
			contextMap["name"] = temp1
			jsonContext := util.ContextSet(*jsonResponse,"90",contextName,contextMap)
			storemessage := "Found " + storename + " near your location, setting " +storename+" as your store"
			jsonResponse.Set(storemessage,"fulfillmentText")

			context.JSON(http.StatusOK, jsonContext.Data())
			return
		case "spot.promotion":
			jsonResponse := gabs.New()
			jsonResponse.Set("Promotion Data", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.order":
			orderResponse := colorlizard.GetOrder()
			jsonResponse := gabs.New()
			jsonResponse.Set(orderResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.parking":
			parkingResponse := colorlizard.GetParking()
			jsonResponse := gabs.New()
			jsonResponse.Set(parkingResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.offers":
			offersResponse := colorlizard.Getoffers()
			jsonResponse := gabs.New()
			jsonResponse.Set(offersResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.payments":
			paymentsResponse := colorlizard.GetPayments()
			jsonResponse := gabs.New()
			jsonResponse.Set(paymentsResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.pets":
			petsResponse := colorlizard.GetPets()
			jsonResponse := gabs.New()
			jsonResponse.Set(petsResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.pharmacy":
			_,contextMap := util.ContextGet(*jsonParsed)
			pharmacymsg := store_details.GetPharmacy(contextMap["store"].Data().(string))
			jsonResponse := gabs.New()
			jsonResponse.Set(pharmacymsg,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		default:
			jsonResponse := gabs.New()
			jsonResponse.Set("Default Response from Webhook", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		}
	})

	r.GET("/store", func(context *gin.Context) {
		store := store_details.GetStoreID("Minneapolis")
		context.JSON(200, gin.H{
			"ready": store,
		})

	})

	return r
}
