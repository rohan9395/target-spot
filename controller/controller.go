package controller

import (
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/target-spot/colorlizard"
	"github.com/target-spot/config"
	"github.com/target-spot/item-search"
	"github.com/target-spot/store-details"
	"github.com/target-spot/util"
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
			_, contextMap := util.ContextGet(*jsonParsed)
			pharmacymsg := store_details.GetStoreAddress(contextMap["store"].Data().(string))
			jsonResponse := gabs.New()
			jsonResponse.Set(pharmacymsg, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.available":
			contextName, searchTermMap := util.ContextGet(*jsonParsed)
			if searchTermMap != nil {
				tcinString, price, title := item_search.GetItemDetails(searchTermMap["itemName.original"].Data().(string))
				jsonResponse := gabs.New()
				temp := gabs.New()
				temp.Set(tcinString)
				temp1 := gabs.New()
				temp1.Set(price)
				temp2 := gabs.New()
				temp2.Set(title)

				searchTermMap["tcin"] = temp
				searchTermMap["itemPrice"] = temp1
				searchTermMap["itemName.original"] = temp2

				jsonContext := util.ContextSet(*jsonResponse, "900", contextName, searchTermMap)
				isItemAvailable := item_search.GetItemAvailability(tcinString, searchTermMap["store"].Data().(string))
				if isItemAvailable {
					jsonResponse.Set(title+" is available at your nearest store "+searchTermMap["name"].Data().(string), "fulfillmentText")
				} else {
					jsonResponse.Set("Get "+title+" on Target.com", "fulfillmentText")
				}
				context.JSON(http.StatusOK, jsonContext.Data())

			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())

			}
			return

		case "spot.getItemPrice":
			_, searchTermMap := util.ContextGet(*jsonParsed)
			if searchTermMap != nil {
				tcinString, price, title := item_search.GetItemDetails(searchTermMap["itemName.original"].Data().(string))
				jsonResponse := gabs.New()
				temp := gabs.New()
				temp.Set(tcinString)
				temp1 := gabs.New()
				temp1.Set(price)
				temp2 := gabs.New()
				temp2.Set(title)

				searchTermMap["tcin"] = temp
				searchTermMap["itemPrice"] = temp1
				searchTermMap["itemName.original"] = temp2

				jsonResponse.Set("Price of "+title+" is "+searchTermMap["itemPrice"].Data().(string), "fulfillmentText")

				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return

		case "spot.price":
			_, searchTermMap := util.ContextGet(*jsonParsed)
			if searchTermMap != nil {
				jsonResponse := gabs.New()
				if searchTermMap["itemName.original"].Data() == nil || searchTermMap["itemName.original"].Data().(string) == "" {
					jsonResponse.Set("No Item to get Price For", "fulfillmentText")
				} else {
					jsonResponse.Set(searchTermMap["itemPrice"].Data().(string), "fulfillmentText")
				}
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return

		case "spot.setstore":
			contextName, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				jsonResponse := gabs.New()

				if contextMap["geo-city"].Data() == nil || contextMap["geo-city"].Data().(string) == "" {
					jsonResponse.Set("Did not get that. Can you please try again ?", "fulfillmentText")
					context.JSON(http.StatusOK, jsonResponse.Data())
				} else {

					store := store_details.GetStoreID(contextMap["geo-city"].Data().(string))
					storename := store_details.GetStoreName(store)

					temp := gabs.New()
					temp.Set(store)
					temp1 := gabs.New()
					temp1.Set(storename)

					contextMap["store"] = temp
					contextMap["name"] = temp1
					jsonContext := util.ContextSet(*jsonResponse, "900", contextName, contextMap)
					storemessage := "Found " + storename + " store near your location, setting " + storename + " as your store"
					jsonResponse.Set(storemessage, "fulfillmentText")

					context.JSON(http.StatusOK, jsonContext.Data())
				}
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.promotion":
			promoResponse := colorlizard.GetPromo()
			jsonResponse := gabs.New()
			jsonResponse.Set(promoResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.order":
			orderResponse := colorlizard.GetOrder()
			jsonResponse := gabs.New()
			jsonResponse.Set(orderResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.parking":
			parkingResponse := colorlizard.GetParking()
			jsonResponse := gabs.New()
			jsonResponse.Set(parkingResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.offers":
			offersResponse := colorlizard.Getoffers()
			jsonResponse := gabs.New()
			jsonResponse.Set(offersResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.payments":
			paymentsResponse := colorlizard.GetPayments()
			jsonResponse := gabs.New()
			jsonResponse.Set(paymentsResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.pets":
			petsResponse := colorlizard.GetPets()
			jsonResponse := gabs.New()
			jsonResponse.Set(petsResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.pharmacy":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetPharmacy(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {

			}
			return
		case "spot.starbucks":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetStarbucks(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.fresh":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetFresh(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.photolab":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetPhotoLab(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.phone":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetStorePhone(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.address":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetStoreAddress(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.timing":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetStoreTiming(contextMap["store"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.id":
			_, contextMap := util.ContextGet(*jsonParsed)
			if contextMap != nil {
				pharmacymsg := store_details.GetStoreAddress(contextMap["number.original"].Data().(string))
				jsonResponse := gabs.New()
				jsonResponse.Set(pharmacymsg, "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())
			}
			return
		case "spot.viewCart":
			contextName, _ := util.ContextGet(*jsonParsed)
			cartResponse := colorlizard.ViewCart(contextName)
			jsonResponse := gabs.New()
			jsonResponse.Set(cartResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		case "spot.addCart":
			contextName, searchTermMap := util.ContextGet(*jsonParsed)
			if searchTermMap != nil {
				jsonResponse := gabs.New()
				_, ok := searchTermMap["itemName.original"]
				if ok {
					itemName := searchTermMap["itemName.original"].Data().(string)
					cartResponse := colorlizard.AddCart(contextName, itemName)
					jsonResponse.Set(cartResponse, "fulfillmentText")
				} else {
					jsonResponse.Set("I dont know which item i should add to cart", "fulfillmentText")
				}
				context.JSON(http.StatusOK, jsonResponse.Data())
			} else {
				jsonResponse := gabs.New()
				jsonResponse.Set("We need your city to get things started, what's your current city?", "fulfillmentText")
				context.JSON(http.StatusOK, jsonResponse.Data())

			}
			return
		case "spot.checkoutCart":
			contextName, _ := util.ContextGet(*jsonParsed)
			cartResponse := colorlizard.CheckoutCart(contextName)
			jsonResponse := gabs.New()
			jsonResponse.Set(cartResponse, "fulfillmentText")
			context.JSON(http.StatusOK, jsonResponse.Data())
			return
		default:
			jsonResponse := gabs.New()
			jsonResponse.Set("I didn't get what you just said, to ge things started what's your current city?", "fulfillmentText")
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
