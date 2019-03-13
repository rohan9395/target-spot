package colorlizard

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
)

func Getoffers() string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, _ := http.Get("https://spot-assist.herokuapp.com/colorlizard/offers")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("offers").Data().(string)
	}
	return ""
}

func GetOrder() string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, _ := http.Get("https://spot-assist.herokuapp.com/colorlizard/orderStatus")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("order_status").Data().(string)
	}
	return ""
}

func GetPayments() string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, _ := http.Get("https://spot-assist.herokuapp.com/colorlizard/payments")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("payments").Data().(string)
	}
	return ""
}

func GetPets() string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, _ := http.Get("https://spot-assist.herokuapp.com/colorlizard/pets")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("pets").Data().(string)
	}
	return ""
}

func GetParking() string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, _ := http.Get("https://spot-assist.herokuapp.com/colorlizard/parking")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("parking").Data().(string)
	}
	return ""
}

func GetPromo() string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, _ := http.Get("https://spot-assist.herokuapp.com/colorlizard/myDummyPromotion")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		// S is shorthand for Search
		return jsonParsed.Path("results").Index(0).Path("promo_desc").Data().(string) + ", " +
			jsonParsed.Path("results").Index(1).Path("promo_desc").Data().(string) + ", " +
			jsonParsed.Path("results").Index(2).Path("promo_desc").Data().(string)
	}
	return ""
}
