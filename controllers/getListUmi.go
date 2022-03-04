package controllers

import (
	"beego/conf"
	"encoding/json"
	"fmt"
	"strconv"
)

type GetParamList struct {
	Branchid string `json:"branchid"`
	Username string `json:"username"`
}

type GetListResponse struct {
	ResponseCode        string        `json:"responseCode"`
	ResponseDescription string        `json:"responseDescription"`
	Data                []DataListGet `json:"data"`
}

type GetListResponseByUserID struct {
	ResponseCode        string      `json:"responseCode"`
	ResponseDescription string      `json:"responseDescription"`
	Data                DataListGet `json:"data"`
}

type DataListGet struct {
	Nik              string `json:"nik"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phoneNumber"`
	BirthDate        string `json:"burthDate"`
	Gender           string `json:"gender"`
	Education        string `json:"educatiuon"`
	Address          string `json:"address"`
	RtRw             string `json:"rtRw"`
	Village          string `json:"village"`
	District         string `json:"district"`
	Province         string `json:"province"`
	PostalCode       string `json:"postalCode"`
	Profession       string `json:"profession"`
	BussinessType    string `json:"businessType"`
	Branch           string `json:"branch"`
	LoanLimit        string `json:"loanLimit"`
	LoanPeriod       string `json:"loanPeriod"`
	IdempotencyKey   string `json:"idempotencyKey"`
	PartnerProductId string `json:"partnerProductId"`
	TransactionID    string `json:"transactionID"`
	CreatedDate      string `json:"CreatedDate"`
	PickByNIP        string `json:"PickByNIP"`
	PickByNama       string `json:"PickByNama"`
	PickByHP         string `json:"PickByHP"`
	PickByDate       string `json:"PickByDate"`
	PickByBranchID   string `json:"PickByBranchID"`
	PickByBranchName string `json:"PickByBranchName"`
	StatusTransaksi  string `json:"StatusTransaksi"`
	UpdatedDate      string `json:"UpdatedDate"`
}

func (c *MainController) GetListNasabah() GetListResponse {

	decoder := json.NewDecoder(c.Ctx.Request.Body)

	dataget := []DataListGet{}
	var dt GetParamList
	err := decoder.Decode(&dt)

	if err != nil {
		fmt.Println(err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = GetListResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad Request",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	fmt.Println("heyyy")

	rows, err := conf.Db.Query(`EXEC GET_LIST_UMICORNER @username = '` + dt.Username + `', @branchid = '` + dt.Branchid + `'`)

	if err != nil {
		fmt.Println(err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = GetListResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad Request",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	defer rows.Close()
	var each = DataListGet{}
	for rows.Next() {

		err = rows.Scan(&each.Nik, &each.Name, &each.PhoneNumber, &each.BirthDate, &each.Gender, &each.Education, &each.Address, &each.RtRw, &each.Village, &each.District, &each.Province, &each.PostalCode, &each.Profession, &each.BussinessType, &each.Branch, &each.LoanLimit, &each.LoanPeriod, &each.IdempotencyKey, &each.PartnerProductId, &each.TransactionID, &each.CreatedDate, &each.PickByNIP, &each.PickByNama, &each.PickByHP, &each.PickByDate, &each.PickByBranchID, &each.PickByBranchName, &each.StatusTransaksi, &each.UpdatedDate)
		if err != nil {
			fmt.Println(err)
			c.Ctx.ResponseWriter.WriteHeader(400)
			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			var response = GetListResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "StatusBadRequest",
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}
		dataget = append(dataget, each)
	}

	// fmt.Println(dataget)

	c.Ctx.ResponseWriter.WriteHeader(200)
	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	var response = GetListResponse{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "Success",
		Data:                dataget,
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response
}

func (c *MainController) GetListNasabahById() GetListResponseByUserID {
	Response := c.Ctx.Input.Param(":transaction_id")

	rows, err := conf.Db.Query(`  SELECT
		nik,
		name,
		phoneNumber,
		birthDate,
		gender,
		education,
		address,
		rtRw,
		village,
		district,
		province,
		postalCode,
		profession,
		bussinessType,
		branch,
		loanLimit,
		loanPeriod,
		idempotencyKey,
		partnerProductId,
		transactionID,
		CreatedDate,
		CASE WHEN PickByNIP IS NULL THEN '-' ELSE PickByNIP END AS PickByNIP,
		CASE WHEN PickByNama IS NULL THEN '-' ELSE PickByNama END  AS PickByNama,
		CASE WHEN PickByHP IS NULL THEN '-' ELSE PickByHP END AS PickByHP,
		CASE WHEN PickByDate IS NULL THEN GETDATE() ELSE PickByDate END AS PickByDate,
		CASE WHEN PickByBranchID IS NULL THEN '-' ELSE PickByBranchID END AS PickByBranchID,
		CASE WHEN PickByBranchName IS NULL THEN '-' ELSE PickByBranchName END AS PickByBranchName,
		CASE WHEN StatusTransaksi IS NULL THEN 'BEING_REGISTERED' ELSE StatusTransaksi END AS StatusTransaksi,
		UpdatedDate
  	FROM UMI_Referal
	WHERE transactionID = '` + Response + `'`)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		fmt.Println("this")
		fmt.Println(err)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = GetListResponseByUserID{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Bad Request",
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	defer rows.Close()
	var each = DataListGet{}
	for rows.Next() {
		err = rows.Scan(&each.Nik, &each.Name, &each.PhoneNumber, &each.BirthDate, &each.Gender, &each.Education, &each.Address, &each.RtRw, &each.Village, &each.District, &each.Province, &each.PostalCode, &each.Profession, &each.BussinessType, &each.Branch, &each.LoanLimit, &each.LoanPeriod, &each.IdempotencyKey, &each.PartnerProductId, &each.TransactionID, &each.CreatedDate, &each.PickByNIP, &each.PickByNama, &each.PickByHP, &each.PickByDate, &each.PickByBranchID, &each.PickByBranchName, &each.StatusTransaksi, &each.UpdatedDate)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(400)
			fmt.Println(err)
			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			var response = GetListResponseByUserID{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "StatusBadRequest",
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	var response = GetListResponseByUserID{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "Success",
		Data:                each,
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response

	// fmt.Println((Response))
}
