package controller

import (
	"github.com/target-spot/offers"
	"github.com/target-spot/payments"
	"github.com/target-spot/pets"
	"github.com/target-spot/pharmacy"
	"github.com/target-spot/store-parking"
	"net/http"
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/target-spot/config"
	"github.com/target-spot/order-status"
)

func GetRouter(endpointMap map[string]util.Endpoint, ready *bool) (r *gin.Engine) {
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
			jsonResponse := gabs.New()
			jsonResponse.Set("Nearest Store Set", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.promotion":
			jsonResponse := gabs.New()
			jsonResponse.Set("Promotion Data", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.order":
			orderResponse := order_status.GetOrder()
			jsonResponse := gabs.New()
			jsonResponse.Set(orderResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.parking":
			parkingResponse := store_parking.GetParking()
			jsonResponse := gabs.New()
			jsonResponse.Set(parkingResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.pharmacy":
			pharmacyResponse := pharmacy.GetPharmacy()
			jsonResponse := gabs.New()
			jsonResponse.Set(pharmacyResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.offers":
			offersResponse := offers.Getoffers()
			jsonResponse := gabs.New()
			jsonResponse.Set(offersResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.payments":
			paymentsResponse := payments.GetPayments()
			jsonResponse := gabs.New()
			jsonResponse.Set(paymentsResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.pets":
			petsResponse := pets.GetPets()
			jsonResponse := gabs.New()
			jsonResponse.Set(petsResponse,"fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		default:
			jsonResponse := gabs.New()
			jsonResponse.Set("Default Response from Webhook", "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		}
	})



	return r
}
