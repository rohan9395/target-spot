package payments

import (
	"crypto/tls"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"net/http"
)
import "fmt"

func GetPayments()(string){
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp,err := http.Get("https://spot-assist.herokuapp.com/colorlizard/payments");
	fmt.Println(err)
	fmt.Println(resp)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes,_ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("payments").Data().(string);
	}
	return ""
}
