package controllers

import (
	"beego/conf"
	"beego/models"
	"encoding/json"
	"fmt"
	"strconv"
)

// "PKM_mekaar/conf"
// "strings"

func (c *MainController) GetTransaction() models.GetResponse {
	// getparam := json.NewDecoder(c.Ctx.Request.Body)
	decoder := json.NewDecoder(c.Ctx.Request.Body)

	authorizationHeader := c.Ctx.Request.Header.Get("Apikey")

	auth := models.Apikey == authorizationHeader

	if auth == false {
		c.Ctx.ResponseWriter.WriteHeader(401)

		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.GetResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Unauthorized",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	dataget := []models.DataGet{}
	var dt models.TransactionGet
	err := decoder.Decode(&dt)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		fmt.Printf("decoder")
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.GetResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Internal server error",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	rows, err := conf.Db.Query(`  SELECT 
	CASE
	  WHEN StatusTransaksi IS NULL THEN 'IN_PROGRESS'
	  WHEN StatusTransaksi = '1' THEN 'DONE_SUCCESS'
	  WHEN StatusTransaksi = '2' THEN 'DONE_FAILED'
	END AS Status, 
	transactionID as TransactionId, 
	CASE WHEN PickByNIP IS NULL THEN 'null' END AS ExecutorId, 
	CASE WHEN PickByNama IS NULL THEN 'null' END AS ExecutorName, 
	CASE WHEN PickByHP IS NULL THEN 'null' END AS ExecutorPhoneNumber, 
	partnerProductId as ExecutionBusinessUnitId, 
	bussinessType as ExecutionBusinessUnitName, 
	loanLimit as LoanRealization
	FROM UMI_Referal
	WHERE transactionID = '` + dt.TransactionId + `'`)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		fmt.Printf("query")
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = models.GetResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Internal server error",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	defer rows.Close()

	for rows.Next() {
		var each = models.DataGet{}
		err = rows.Scan(&each.Status, &each.TransactionId, &each.ExecutorId, &each.ExecutorName, &each.ExecutorPhoneNumber, &each.ExecutionBusinessUnitId, &each.ExecutionBusinessUnitName, &each.LoanRealization)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(500)
			fmt.Printf(err.Error())
			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			var response = models.GetResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "Internal server error",
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}
		dataget = append(dataget, each)
	}

	c.Ctx.ResponseWriter.WriteHeader(200)

	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	var datatest = models.DataGet{
		Status:                    dataget[0].Status,
		TransactionId:             dataget[0].TransactionId,
		ExecutorId:                dataget[0].ExecutorId,
		ExecutorName:              dataget[0].ExecutorName,
		ExecutorPhoneNumber:       dataget[0].ExecutorPhoneNumber,
		ExecutionBusinessUnitId:   dataget[0].ExecutionBusinessUnitId,
		ExecutionBusinessUnitName: dataget[0].ExecutionBusinessUnitName,
		LoanRealization:           dataget[0].LoanRealization,
	}

	var response = models.GetResponse{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "success",
		Data:                datatest,
	}

	c.Data["json"] = response

	c.ServeJSON()

	return response

	// var param models.TransactionGet

	// reqBody, err := ioutil.ReadAll(c.Ctx.Request.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", reqBody)

}
