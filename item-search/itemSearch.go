package item_search

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/rs/zerolog/log"
)

type HTTPClient struct {
	client *http.Client
	pool   *x509.CertPool
}

const itemSearchURL = "https://redsky.target.com/v1/plp/search/?count=24&offset=0&keyword=%s"
const itemAvailabilityURL = "https://redsky.target.com/v1/location_details/%s?storeId=%s"

func GetItemDetails(searchTerm string) (string, string, string) {

	jsonParsed1 := makeItemSearchDetails(searchTerm)

	data1, err := jsonParsed1.Path("search_response.items.Item").Children()
	if err != nil {
		//no data for item found . Return
		return "", "", ""
	} else {
		for _, child1 := range data1 {
			tcinString := ""
			title := ""
			itemPrice := ""
			if child1.Exists("tcin") {
				tcinString = child1.Path("tcin").Data().(string)
			}
			if child1.Exists("title") {
				title = child1.Path("title").Data().(string)
			}
			if child1.Exists("list_price.formatted_price") {
				itemPrice = child1.Path("list_price.formatted_price").Data().(string)
			} else if child1.Exists("offer_price.formatted_price") {
				itemPrice = child1.Path("offer_price.formatted_price").Data().(string)
			} else {
				itemPrice = "is not available"
			}
			return tcinString, itemPrice, title
		}
	}
	return "", "", ""

}

func GetItemAvailability(searchTcin string, store string) bool {

	jsonParsed1 := makeItemAvailabilityBool(searchTcin, store)

	data1, err := jsonParsed1.Path("product.available_to_promise_store.products").Children()
	if err != nil {
		//did not get a response from redsky. Return
		return false
	} else {
		for _, child := range data1 {
			data2, err := child.Path("locations").Children()
			if err != nil {
				//did not get a response from redsky. Return
				return false
			} else {
				for _, child2 := range data2 {
					isItemAvailable := child2.Path("onhand_quantity").Data().(float64)
					if isItemAvailable > 0.0 {
						return true
					}
				}
			}
		}
	}
	return false
}

func makeItemSearchDetails(searchTerm string) *gabs.Container {
	h := HTTPClient{}
	url := fmt.Sprintf(itemSearchURL, strings.Replace(searchTerm, " ", "+", -1))

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/json")
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new http request.")
	}

	h.client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	response, err := h.client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Error making call to redsky api")
	}

	responseData, err := ioutil.ReadAll(response.Body)

	jsonParsed, err := gabs.ParseJSON(responseData)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing response of redsky api")
	}
	return jsonParsed

}

func makeItemAvailabilityBool(searchTerm string, store string) *gabs.Container {
	h := HTTPClient{}
	url := fmt.Sprintf(itemAvailabilityURL, searchTerm, store)

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/json")
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new http request.")
	}

	h.client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	response, err := h.client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Error making call to redsky api")
	}

	responseData, err := ioutil.ReadAll(response.Body)

	jsonParsed, err := gabs.ParseJSON(responseData)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing response of redsky api")
	}
	return jsonParsed

}
