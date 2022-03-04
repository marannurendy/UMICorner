package controllers

import (
	"beego/conf"
	"beego/models"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func genHMAC256(ciphertext, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(ciphertext))
	hmac := mac.Sum(nil)
	return hmac
}

func GetToken() string {
	endpoint := "https://partner.api.bri.co.id/oauth/client_credential/accesstoken?grant_type=client_credentials"
	data := url.Values{}
	data.Set("client_id", "051tUdGPdW9SCBwImUmYc52UYKsG1rGo")
	data.Set("client_secret", "YGGW7BsilqPGmOAd")

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload

	if err != nil {
		fmt.Println(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	// log.Println(res.Status)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	// log.Println(string(body))

	var dt = (body)

	var dataShit WebAuth

	err = json.Unmarshal([]byte(dt), &dataShit)

	return dataShit.Access_token
}

func GetSignature(token string, method string, path string, body string, timestamp string) string {

	secretKey := "YGGW7BsilqPGmOAd"
	payload := "path=" + path + "&verb=" + method + "&token=Bearer " + token + "&timestamp=" + timestamp + "&body=" + body
	//payload := `path=/v1.0/umi/webauth&verb=POST&token=Bearer 7U8rVXKUQVyjCnMnXsBzILu0Mxnt&timestamp=2021-07-30T17:04:05.000Z&body={"sellerId":"00219333","name":"TOYIB","businessUnit":"Bandung","businessUnitId":"0206","title":"Mantri"}`

	fmt.Println(payload)
	hmac := genHMAC256([]byte(payload), []byte(secretKey))
	stringHmac := b64.StdEncoding.EncodeToString(hmac)
	fmt.Println(stringHmac)
	return stringHmac
}

func (c *MainController) UpdateTransaction() models.TransactionUpdateResponse {
	decoder := json.NewDecoder(c.Ctx.Request.Body)

	authorizationHeader := c.Ctx.Request.Header.Get("Apikey")

	auth := models.Apikey == authorizationHeader

	if auth == false {
		c.Ctx.ResponseWriter.WriteHeader(401)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionUpdateResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Unauthorized",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	var dt models.TransactionUpdate
	err := decoder.Decode(&dt)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)

		fmt.Println("decoder" + err.Error())

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionUpdateResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Internal server error",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	ctx := context.Background()
	tx, err := conf.Db.BeginTx(ctx, nil)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(404)

		fmt.Println("context" + err.Error())

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionUpdateResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Resources Not Found",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	// _, err = tx.ExecContext(ctx, `UPDATE UMI_Referal
	// SET StatusTransaksi='`+dt.Status+`',
	// 	PickByNIP='`+dt.ExecutorId+`',
	// 	PickByNama='`+dt.ExecutorName+`',
	// 	PickByHP='`+dt.ExecutorPhoneNumber+`',
	// 	PickByBranchID='`+dt.ExecutionBusinessUnitId+`',
	// 	PickByBranchName='`+dt.ExecutionBusinessUnitName+`',
	// 	PickByDate=CAST(GETDATE() as DATETIME)
	// 	where TransactionId='`+dt.TransactionId+`'`)

	_, err = tx.ExecContext(ctx, `UPDATE UMI_Referal
	SET StatusTransaksi='`+dt.Status+`',
		PickByNama='`+dt.ExecutorName+`',
		PickByBranchID='`+dt.ExecutionBusinessUnitId+`',
		PickByBranchName='`+dt.ExecutionBusinessUnitName+`',
		PickByDate=CAST(GETDATE() as DATETIME)
		where TransactionId='`+dt.TransactionId+`'`)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(404)

		fmt.Println("context" + err.Error())

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionUpdateResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Resources Not Found",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	err = tx.Commit()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		fmt.Println("commit" + err.Error())

		var response = models.TransactionUpdateResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad request or invalid input validation",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	} else {
		// TEST
		token := GetToken()
		method := "PATCH"
		path := "/v1.0/umi/referral"
		timestamp := (time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
		body := []byte(`{"status": "` + dt.Status + `","transactionId": "` + dt.TransactionId + `","executorId": "` + dt.ExecutorId + `","executorName": "` + dt.ExecutorName + `","executorPhoneNumber": "` + dt.ExecutorPhoneNumber + `","executionBusinessUnitId": "` + dt.ExecutionBusinessUnitId + `","executionBusinessUnitName": "` + dt.ExecutionBusinessUnitId + `","loanRealization": "2000000"}`)
		signature := GetSignature(token, method, path, string(body), timestamp)

		lastPoint := "https://partner.api.bri.co.id/v1.0/umi/referral"

		Otherclient := &http.Client{}
		s, err := http.NewRequest(method, lastPoint, bytes.NewBuffer(body)) // URL-encoded payload

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(400)

			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			// fmt.Println("commit" + err.Error())

			var response = models.TransactionUpdateResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "Bad request or invalid input validation",
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}

		var bearer = "Bearer " + token

		s.Header.Add("Content-Type", "application/json")
		s.Header.Add("BRI-Signature", signature)
		s.Header.Add("BRI-Timestamp", timestamp)
		s.Header.Add("Authorization", bearer)

		resul, err := Otherclient.Do(s)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(400)

			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			// fmt.Println("commit" + err.Error())

			var response = models.TransactionUpdateResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "Bad request or invalid input validation",
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}

		// log.Println(resul.Status)
		defer resul.Body.Close()

		_, err = ioutil.ReadAll(resul.Body)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(400)

			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			// fmt.Println("commit" + err.Error())

			var response = models.TransactionUpdateResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "Bad request or invalid input validation",
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}

		// fmt.Println(string(bodyOther))
	}

	c.Ctx.ResponseWriter.WriteHeader(201)

	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	var response = models.TransactionUpdateResponse{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "created",
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response

}

func (c *MainController) PostTransaction() models.TransactionResponse {
	decoder := json.NewDecoder(c.Ctx.Request.Body)
	var transactionID string

	authorizationHeader := c.Ctx.Request.Header.Get("Apikey")

	auth := models.Apikey == authorizationHeader

	if auth == false {
		c.Ctx.ResponseWriter.WriteHeader(401)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var dataid = models.DetData{
			Transactionid: transactionID,
		}

		var response = models.TransactionResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Unauthorized",
			Data:                dataid,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	var dt models.TransactionPost
	err := decoder.Decode(&dt)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var dataid = models.DetData{
			Transactionid: transactionID,
		}

		var response = models.TransactionResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Internal server error",
			Data:                dataid,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	ctx := context.Background()
	tx, err := conf.Db.BeginTx(ctx, nil)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(404)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var dataid = models.DetData{
			Transactionid: transactionID,
		}

		var response = models.TransactionResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Resources Not Found",
			Data:                dataid,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	insertTransac := `INSERT INTO UMI_Referal
	OUTPUT INSERTED.transactionID
	VALUES (
		'` + dt.Nik + `',
		'` + dt.Name + `',
		'` + dt.PhoneNumber + `',
		'` + dt.Birthdate + `',
		'` + dt.Gender + `',
		'` + dt.Education + `',
		'` + dt.Address + `',
		'` + dt.RtRw + `',
		'` + dt.Village + `',
		'` + dt.District + `',
		'` + dt.Province + `',
		'` + dt.PostalCode + `',
		'` + dt.Profession + `',
		'` + dt.BusinessType + `',
		'` + dt.Branch + `',
		'` + dt.LoanLimit + `',
		'` + dt.LoanPeriod + `',
		'` + dt.IdempotencyKey + `',
		'` + dt.PartnerProductId + `',
		NEWID(),
		CAST(GETDATE() as DATETIME),
		NULL,
		NULL,
		NULL,
		NULL,
		NULL,
		NULL,
		NULL,
		CAST(GETDATE() as DATETIME)
	)`

	stmt, err := tx.Prepare(insertTransac)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var dataid = models.DetData{
			Transactionid: transactionID,
		}

		var response = models.TransactionResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad request or invalid input validation",
			Data:                dataid,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	defer stmt.Close()

	err = stmt.QueryRow().Scan(&transactionID)

	if err != nil {
		// Incase we find any error in the query execution, rollback the transaction
		c.Ctx.ResponseWriter.WriteHeader(400)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var dataid = models.DetData{
			Transactionid: transactionID,
		}

		var response = models.TransactionResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad request or invalid input validation",
			Data:                dataid,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	err = tx.Commit()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var dataid = models.DetData{
			Transactionid: transactionID,
		}

		var response = models.TransactionResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad request or invalid input validation",
			Data:                dataid,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	c.Ctx.ResponseWriter.WriteHeader(201)

	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	var dataid = models.DetData{
		Transactionid: transactionID,
	}

	var response = models.TransactionResponse{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "created",
		Data:                dataid,
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response

}

func (c *MainController) PostInquiry() models.TransactionIncResponse {
	decoder := json.NewDecoder(c.Ctx.Request.Body)

	authorizationHeader := c.Ctx.Request.Header.Get("Apikey")

	auth := models.Apikey == authorizationHeader

	if auth == false {
		c.Ctx.ResponseWriter.WriteHeader(401)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionIncResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Unauthorized",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	var dt models.InquiryPost
	err := decoder.Decode(&dt)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionIncResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Internal server error",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	insertTransac := `SELECT 
	CASE
	  WHEN StatusTransaksi IS NULL THEN 'IN_PROGRESS'
	  WHEN StatusTransaksi = '1' THEN 'DONE_SUCCESS'
	  WHEN StatusTransaksi = '2' THEN 'DONE_FAILED'
	  ELSE StatusTransaksi
	END AS Status, 
	transactionID as TransactionId, 
	CASE WHEN PickByNIP IS NULL THEN 'null' ELSE PickByNIP END AS ExecutorId, 
	CASE WHEN PickByNama IS NULL THEN 'null' ELSE PickByNama END AS ExecutorName, 
	CASE WHEN PickByHP IS NULL THEN 'null' ELSE PickByHP END AS ExecutorPhoneNumber, 
	CASE WHEN PickByBranchID IS NULL THEN 'null' ELSE PickByBranchID END AS ExecutionBusinessUnitId, 
	CASE WHEN PickByBranchName IS NULL THEN 'null' ELSE PickByBranchName END AS ExecutionBusinessUnitName, 
	loanLimit as LoanRealization
	FROM UMI_Referal WITH(NOLOCK)
	WHERE transactionID = '` + dt.TransactionId + `'`

	prep, err := conf.Db.Prepare(insertTransac)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionIncResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad request or invalid input validation",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	defer prep.Close()
	var each = models.DataInquiry{}
	err = prep.QueryRow().Scan(&each.Status, &each.TransactionId, &each.ExecutorId, &each.ExecutorName, &each.ExecutorPhoneNumber, &each.ExecutionBusinessUnitId, &each.ExecutionBusinessUnitName, &each.LoanRealization)
	if err != nil {

		fmt.Println(err.Error())
		c.Ctx.ResponseWriter.WriteHeader(400)
		fmt.Println(err.Error())
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.TransactionIncResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad request or invalid input validation",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	c.Ctx.ResponseWriter.WriteHeader(200)

	// fmt.Println(dt.TransactionId)

	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	var response = models.TransactionIncResponse{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "Success",
		Data:                &each,
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response

}
