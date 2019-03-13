package store_details

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Jeffail/gabs"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
)
type HTTPClient struct {
	client *http.Client
	pool   *x509.CertPool
}

func GetStoreName(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	name := jsonParsed1.Path("Name").String()
	name = strings.Trim(name,"\"")

	name = strings.Trim(name,"\\")
	return "Your default store will be set as " + name
}

func GetStoreID(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	id := jsonParsed1.Path("ID").String()
	id = strings.Trim(id,"\"")

	id = strings.Trim(id,"\\")
	return id
}

func GetFarmacy(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	capability := jsonParsed1.Path("Capability.CapabilityName").String()

	id := strings.Contains(capability,"CVS pharmacy")
	if id{
		return "Yes, CVS pharmacy is available in this store"
	}else{
		return "No, pharmacy is not available in this store"
	}

}

func GetStarbucks(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	capability := jsonParsed1.Path("Capability.CapabilityName").String()

	id := strings.Contains(capability,"Starbucks")
	if id{
		return "Yes, Starbucks is available in this store"
	}else{
		return "No, Starbucks is not available in this store"
	}

}

func GetFresh(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	capability := jsonParsed1.Path("Capability.CapabilityName").String()

	id := strings.Contains(capability,"Fresh Grocery")
	if id{
		return "Yes, Fresh Grocery is available in this store"
	}else{
		return "No, Fresh Grocery is not available in this store"
	}

}

func GetPhotoLab(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	capability := jsonParsed1.Path("Capability.CapabilityName").String()

	id := strings.Contains(capability,"Photo Lab")
	if id{
		return "Yes, Photo Lab is available in this store"
	}else{
		return "No, Photo Lab is not available in this store"
	}

}

func GetStorePhone(city string ) string {

	jsonParsed1 := makeStoreDetailsBaseCall(city)
	TelephoneNumber := jsonParsed1.Path("TelephoneNumber")
	phone := TelephoneNumber.Path("PhoneNumber").Index(0).String()

	phone = strings.Trim(phone,"\"")
	phone = strings.Trim(phone,"\\")
	phone = strings.Trim(phone,"\\")

	return "Phone Number of store is " + phone
}

func makeStoreDetailsBaseCall(city string) *gabs.Container{
	h := HTTPClient{}
	url := "https://redsky.target.com/v2/stores/nearby/Minneapolis?locale=en-US&limit=20&range=250"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/json")
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new http request.")
	}

	h.client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: h.pool}}}

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