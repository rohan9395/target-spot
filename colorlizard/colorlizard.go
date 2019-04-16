package colorlizard

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
)

var carts = make(map[string]string)

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

func ViewCart(userName string) string {
	_, ok := carts[userName]
	if ok {
		return carts[userName]
	} else {
		return "No items in cart"
	}
}

func AddCart(userName string, item string) string {
	_, ok := carts[userName]
	if ok {
		carts[userName] += ", " + item
	} else {
		carts[userName] = "Your cart has " + item
	}
	return "item added to cart"
}

func CheckoutCart(userName string) string {
	_, ok := carts[userName]
	if ok {
		delete(carts, userName)
		return "Your cart has been successfully checked out."
	} else {
		return "No items to check out"
	}

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
