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
	tcinString := jsonParsed1.Path("search_response.items.Item").Index(0).Path("tcin").Data().(string)
	title := jsonParsed1.Path("search_response.items.Item").Index(0).Path("title").Data().(string)
	itemPrice := jsonParsed1.Path("search_response.items.Item").Index(0).Path("list_price.formatted_price").Data().(string)
	return tcinString, itemPrice, title
}

func GetItemAvailability(searchTcin string, store string) bool {

	jsonParsed1 := makeItemAvailabilityBool(searchTcin, store)
	isItemAvailable := jsonParsed1.Path("product.available_to_promise_store.products").Index(0).Path("locations").Index(0).Path("onhand_quantity").Data().(float64)
	if isItemAvailable > 0.0 {
		return true
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
