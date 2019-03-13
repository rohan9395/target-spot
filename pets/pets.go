package pets

import (
	"crypto/tls"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"net/http"
)
import "fmt"

func GetPets()(string){
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp,err := http.Get("https://spot-assist.herokuapp.com/colorlizard/pets");
	fmt.Println(err)
	fmt.Println(resp)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes,_ := ioutil.ReadAll(resp.Body)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		return jsonParsed.Path("pets").Data().(string);
	}
	return ""
}
