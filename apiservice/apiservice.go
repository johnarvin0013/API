package apiservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/johnarvin0013/fiber/fiberutils"

	"github.com/gofiber/fiber/v2"
)

//kukuhanin
type IPAPIResponse struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
}

//ibabato or vice-versa
type GMapURLResponse struct {
	URL     string `json:"url"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func GetGMapURL(c *fiber.Ctx) error {
	var ipapiResponse IPAPIResponse
	ipAddress := fiberutils.GetParamValue(c, "ip", "No Ip Address! Please enter IP address.")
	requestURL := "http://ip-api.com/json/" + ipAddress
	ipAPIRequest(requestURL, &ipapiResponse)
	googleMapsURL := fmt.Sprintf("https://www.google.com/maps/@%f,%f,15z", ipapiResponse.Latitude, ipapiResponse.Longitude)

	var gMapURLResponse GMapURLResponse
	gMapURLResponse.URL = googleMapsURL
	gMapURLResponse.Country = ipapiResponse.Country
	gMapURLResponse.City = ipapiResponse.City
	c.JSON(gMapURLResponse)
	return nil
}

func ipAPIRequest(requestURL string, ipapiResponse *IPAPIResponse) {
	response, _ := http.Get(requestURL)
	// fmt.Println("Processing")
	// time.Sleep(3 * time.Second)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("Finish")
	json.Unmarshal(body, &ipapiResponse)
}
