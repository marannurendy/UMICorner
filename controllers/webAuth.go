package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// "PKM_mekaar/conf"
// "strings"

type WebAuth struct {
	Refresh_token_expires_in string `json:"refresh_token_expires_in"`
	Api_product_list         string `json:"api_product_list"`
	Api_product_list_json    string `json:"api_product_list_json"`
	Organization_name        string `json:"organization_name"`
	DeveloperEmail           string `json:"developer.email"`
	Token_type               string `json:"token_type"`
	Issued_at                string `json:"issued_at"`
	Client_id                string `json:"client_id"`
	Access_token             string `json:"access_token"`
	Application_name         string `json:"application_name"`
	Scope                    string `json:"scope"`
	Expires_in               string `json:"expires_in"`
	Refresh_count            string `json:"refresh_count"`
	Status                   string `json:"status"`
}

type WebAuthResponse struct {
	ResponseCode        string      `json:"responseCode"`
	ResponseDescription string      `json:"responseDescription"`
	Data                WebAuthData `json:"data"`
}

type WebAuthData struct {
	Expired    string `json:"epired"`
	Token      string `json:"token"`
	WebviewUrl string `json:"webviewUrl"`
}

type GetWeb struct {
	SellerId       string `json:"sellerid"`
	Name           string `json:"name"`
	BusinessUnit   string `json:"businessUnit"`
	BusinessUnitId string `json:"businessUnitId"`
	Title          string `json:"title"`
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *MainController) WebAuth() WebAuthResponse {

	decoder := json.NewDecoder(c.Ctx.Request.Body)
	var fuck GetWeb
	err := decoder.Decode(&fuck)

	lastPoint := "https://partner.api.bri.co.id/v1.0/umi/webauth"

	var Anotherdata = []byte(`{
		"sellerId": "` + fuck.SellerId + `",
		"name": "` + fuck.Name + `",
		"businessUnit": "` + fuck.BusinessUnit + `",
		"businessUnitId": "` + fuck.BusinessUnitId + `",
		"title": "` + fuck.Title + `"
	}`)

	// fmt.Println(`yah yang ini nah {
	// 	"sellerId": "` + fuck.SellerId + `",
	// 	"name": "` + fuck.Name + `",
	// 	"businessUnit": "` + fuck.BusinessUnit + `",
	// 	"businessUnitId": "` + fuck.BusinessUnitId + `",
	// 	"title": "` + fuck.Title + `"
	// }`)

	token := GetToken()
	method := "POST"
	path := "/v1.0/umi/webauth"
	timestamp := (time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
	body := Anotherdata
	signature := GetSignature(token, method, path, string(body), timestamp)

	Otherclient := &http.Client{}
	s, err := http.NewRequest(method, lastPoint, bytes.NewBuffer(Anotherdata)) // URL-encoded payload

	if err != nil {
		fmt.Println(err)
	}

	var bearer = "Bearer " + token

	s.Header.Add("Content-Type", "application/json")
	s.Header.Add("BRI-Signature", signature)
	s.Header.Add("BRI-Timestamp", timestamp)
	s.Header.Add("Authorization", bearer)

	resul, err := Otherclient.Do(s)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("yang ini " + resul.Status)
	defer resul.Body.Close()

	dataList := WebAuthResponse{}

	json.NewDecoder(resul.Body).Decode(&dataList)

	// fmt.Println(dataList.Data)

	// bodyOther, err := ioutil.ReadAll(resul.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// a := string(bodyOther)
	// fmt.Println(a)
	// data := WebAuthResponse{}
	// json.Unmarshal([]byte(a), data)
	// fmt.Println(data.ResponseCode)

	// statcod := c.Ctx.ResponseWriter.Status
	// statusCode := strconv.Itoa(statcod)

	// fmt.Println(dataList)

	var response = WebAuthResponse{
		ResponseCode:        dataList.ResponseCode,
		ResponseDescription: dataList.ResponseDescription,
		Data:                dataList.Data,
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response

}
