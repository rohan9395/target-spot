package store_details

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/rs/zerolog/log"
)

type HTTPClient struct {
	client *http.Client
	pool   *x509.CertPool
}

const storeurl = "https://redsky.target.com/v2/stores/location/%s"
const cityurl = "https://redsky.target.com/v2/stores/nearby/%s?locale=en-US&limit=20&range=250"

func GetStoreName(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	name := jsonParsed1.Index(0).Path("name").Data().(string)
	return name
}

func GetStoreID(city string) string {

	jsonParsed1 := makeStoreDetailsbyCity(city)
	id := jsonParsed1.Path("ID").Data().(float64)
	storeId := fmt.Sprintf("%d", int(id))
	return storeId
}

func GetPharmacy(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	capability := jsonParsed1.Path("capabilities").String()

	id := strings.Contains(capability, "CVS pharmacy")
	if id {
		return "Yes, CVS pharmacy is available in this store"
	} else {
		return "No, pharmacy is not available in this store"
	}

}

func GetStarbucks(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	capability := jsonParsed1.Path("capabilities").String()

	id := strings.Contains(capability, "Starbucks")
	if id {
		return "Yes, Starbucks is available in this store"
	} else {
		return "No, Starbucks is not available in this store"
	}

}

func GetFresh(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	capability := jsonParsed1.Path("capabilities").String()

	id := strings.Contains(capability, "Fresh Grocery")
	if id {
		return "Yes, Fresh Grocery is available in this store"
	} else {
		return "No, Fresh Grocery is not available in this store"
	}

}

func GetPhotoLab(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	capability := jsonParsed1.Path("capabilities").String()

	id := strings.Contains(capability, "Photo Lab")
	if id {
		return "Yes, Photo Lab is available in this store"
	} else {
		return "No, Photo Lab is not available in this store"
	}

}

func GetStorePhone(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	TelephoneNumber := jsonParsed1.Index(0).Path("phoneNumbers")
	phone := TelephoneNumber.Index(0).Path("phoneNumber").Data().(string)

	return "Phone Number of store is " + phone
}

func GetStoreAddress(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	address := jsonParsed1.Index(0).Path("address.formattedAddress").Data().(string)

	return "Stores address is " + address
}

func GetStoreTiming(storeId string) string {

	jsonParsed1 := makeStoreDetails(storeId)
	starttime := jsonParsed1.Index(0).Path("operatingHours").Index(1).Path("timePeriod.beginTime").Data().(string)
	endtime := jsonParsed1.Index(0).Path("operatingHours").Index(1).Path("timePeriod.thruTime").Data().(string)
	start := strings.Split(starttime, ":")
	end := strings.Split(endtime, ":")
	startint, _ := strconv.Atoi(start[0])
	endint, _ := strconv.Atoi(end[0])

	if len(start) > 0 && len(end) > 0 {
		return fmt.Sprintf("Stores opens at %d hours in the morning and closes at %d hours in the evening", startint, endint)
	} else {
		return "I can't find store timings"
	}
}

func makeStoreDetailsbyCity(city string) *gabs.Container {
	h := HTTPClient{}
	url := fmt.Sprintf(cityurl, city)

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
	jsonParsed1 := jsonParsed.Path("Locations.Location").Index(0)
	return jsonParsed1

}

func makeStoreDetails(storeId string) *gabs.Container {
	h := HTTPClient{}
	url := fmt.Sprintf(storeurl, storeId)

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
