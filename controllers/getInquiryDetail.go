package controllers

import (
	"beego/conf"
	"beego/models"
	"fmt"
	"strconv"
	"strings"
)

type GetRecommendationResponse struct {
	ResponseCode        string              `json:"responseCode"`
	ResponseDescription string              `json:"responseDescription"`
	Data                *recommendationData `json:"data"`
}

type recommendationData struct {
	RecommendationId       string                       `json:"recommendationId"`
	ExistingProducs        []dataExistingProducs        `json:"existingProducts"`
	ProductRecommendations []dataProductRecommendations `json:"productRecommendations"`
	Name                   string                       `json:"name"`
	IdentityNumber         string                       `json:"identityNumber"`
	PhoneNumber            dataPhoneNumber              `json:"phoneNumber"`
	Gender                 string                       `json:"gender"`
	Age                    string                       `json:"age"`
	IdentityAddress        dataIdentityAddress          `json:"identityAddress"`
	DomicileAddress        dataDomicileAddress          `json:"domicileAddress"`
	Education              string                       `json:"education"`
	AccountType            string                       `json:"accountType"`
	WorkUnitCode           string                       `json:"workUnitCode"`
	BusinessType           *dataBusinessType            `json:"businessType"`
	Branchid               string                       `json:"branchId"`
	KelompokID             string                       `json:"kelompokID"`
	Entity                 string                       `json:"entity"`
	Type                   string                       `json:"type"`
	EmployeeID             string                       `json:"employeeID"`
	Score                  string                       `json:"score"`
}

type dataExistingProducs struct {
	PartnerProductId   string `json:"partnerProductId"`
	PartnerProductName string `json:"partnerProductName"`
	ProductType        string `json:"productType"`
	Plafond            string `json:"plafond"`
	BakiDebet          string `json:"bakiDebet"`
	Angsuran           string `json:"angsuran"`
	DueDate            string `json:"dueDate"`
}

type dataProductRecommendations struct {
	Id            string `json:"id"`
	ProjectEntity string `json:"projectEntity"`
}

type dataPhoneNumber struct {
	Mobile string `json:"mobile"`
	Home   string `json:"home"`
}

type dataIdentityAddress struct {
	Name       string `json:"name"`
	Rtrw       string `json:"rtrw"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Village    string `json:"village"`
	PostalCode string `json:"postalCode"`
}

type dataDomicileAddress struct {
	Name       string `json:"name"`
	Rtrw       string `json:"rtrw"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Village    string `json:"village"`
	PostalCode string `json:"postalCode"`
}

type dataBusinessType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type arrayDataProductRecommendation struct {
	Id            string `json:"id"`
	ProjectEntity string `json:"projectEntity"`
}

func (c *MainController) GetRecommendationId() GetRecommendationResponse {
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	Response := c.Ctx.Input.Param(":recommendationId")
	authorizationHeader := c.Ctx.Request.Header.Get("Apikey")

	auth := models.Apikey == authorizationHeader

	if auth == false {
		c.Ctx.ResponseWriter.WriteHeader(401)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)
		var response = GetRecommendationResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "Unauthorized",
			Data:                nil,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return response
	}

	dataProductExisting := []dataExistingProducs{}
	rows, err := conf.Db2.Query(`SELECT
		partnerProductId,
		partnerProductName,
		productType,
		plafond,
		bakiDebit,
		angsuran,
		dueDate,
		productRecommendations,
		projectEntity
		FROM Z_KBUMN_FINAL_API_Inquiry_Detail WHERE RecommendationID = '` + Response + `'`)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		fmt.Println("this")
		fmt.Println(err)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = GetRecommendationResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "bad request or invalid input validation",
			Data:                nil,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return response
	}

	defer rows.Close()
	var existingProducts = dataExistingProducs{}
	var arrayProductRecommendation = arrayDataProductRecommendation{}
	for rows.Next() {
		err = rows.Scan(&existingProducts.PartnerProductId, &existingProducts.PartnerProductName, &existingProducts.ProductType, &existingProducts.Plafond, &existingProducts.BakiDebet, &existingProducts.Angsuran, &existingProducts.DueDate, &arrayProductRecommendation.Id, &arrayProductRecommendation.ProjectEntity)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(404)
			fmt.Println(err)
			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			var response = GetRecommendationResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "resources not found",
				Data:                nil,
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}
		// dataRecommendations = append(dataRecommendations, productRecommendation)
		dataProductExisting = append(dataProductExisting, existingProducts)
	}

	var stringRecoId = arrayProductRecommendation.Id
	var stringProEnt = arrayProductRecommendation.ProjectEntity

	idReco := strings.Split(stringRecoId, ",")
	proEnt := strings.Split(stringProEnt, ",")

	dataRecom := []dataProductRecommendations{}
	var arrayRecom = arrayDataProductRecommendation{}

	for i := 0; i < len(idReco); i++ {
		arrayRecom.Id = idReco[i]
		arrayRecom.ProjectEntity = proEnt[i]

		fmt.Println(idReco[i])

		if idReco[i] != "-" {
			dataRecom = append(dataRecom, dataProductRecommendations(arrayRecom))
		}
	}

	// fmt.Println(dataRecom)

	var dataRec = recommendationData{
		ProductRecommendations: dataRecom,
		ExistingProducs:        dataProductExisting,
	}

	// START Get Data Identity Address

	rows2, err := conf.Db2.Query(`SELECT

		RecommendationID,
		Name,
		identityNumber,
		gender,
		age,
		education,
		CASE WHEN branchID IS NULL THEN '-' ELSE branchID END AS branchID,
		entity,
		handPhoneNumber,
		homePhoneNumber,
		nameIdentity,
		rtrwIdentity,
		provinceIdentity,
		cityIdentity,
		districtIdentity,
		villageIdentity,
		postalCodeIdentity,
		nameDomicile,
		rtrwDomicile,
		provincedomicile,
		citydomicile,
		districtdomicile,
		villageDomicile,
		postalCodeDomicile,
		accountType,
		businessTypeID,
		businessTypeName,
		type,
		employeeID,
		score,
		kelompokID,
		namaKelompok
		FROM Z_KBUMN_FINAL_API_Inquiry_Detail WHERE RecommendationID = '` + Response + `'`)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		fmt.Println("this")
		fmt.Println(err)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = GetRecommendationResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "bad request or invalid input validation",
			Data:                nil,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	defer rows2.Close()
	var data = dataRec
	var dataBusType = dataBusinessType{}

	var datas []recommendationData

	for rows2.Next() {
		err = rows2.Scan(
			&data.RecommendationId,
			&data.Name,
			&data.IdentityNumber,
			&data.Gender,
			&data.Age,
			&data.Education,
			&data.Branchid,
			&data.Entity,
			&data.PhoneNumber.Mobile,
			&data.PhoneNumber.Home,
			&data.IdentityAddress.Name,
			&data.IdentityAddress.Rtrw,
			&data.IdentityAddress.Province,
			&data.IdentityAddress.City,
			&data.IdentityAddress.District,
			&data.IdentityAddress.Village,
			&data.IdentityAddress.PostalCode,
			&data.DomicileAddress.Name,
			&data.DomicileAddress.Rtrw,
			&data.DomicileAddress.Province,
			&data.DomicileAddress.City,
			&data.DomicileAddress.District,
			&data.DomicileAddress.Village,
			&data.DomicileAddress.PostalCode,
			&data.AccountType,
			&dataBusType.Id,
			&dataBusType.Name,
			&data.Type,
			&data.EmployeeID,
			&data.Score,
			&data.KelompokID,
			&data.WorkUnitCode)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(400)
			fmt.Println(err)
			statcod := c.Ctx.ResponseWriter.Status
			statusCode := strconv.Itoa(statcod)

			var response = GetRecommendationResponse{
				ResponseCode:        "0" + statusCode,
				ResponseDescription: "StatusBadRequest",
				Data:                nil,
			}
			c.Data["json"] = response

			c.ServeJSON()

			return response
		}

		datas = append(datas, data)
	}

	rows2.Close()

	if len(datas) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(404)
		fmt.Println(err)
		statcod := c.Ctx.ResponseWriter.Status
		statusCode := strconv.Itoa(statcod)

		var response = GetRecommendationResponse{
			ResponseCode:        "0" + statusCode,
			ResponseDescription: "resources not found",
			Data:                nil,
		}
		c.Data["json"] = response

		c.ServeJSON()

		return response
	}

	data.BusinessType = &dataBusType
	// END Data Identity Address

	// fmt.Println(each)

	c.Ctx.ResponseWriter.WriteHeader(200)
	fmt.Println("this")
	fmt.Println(err)
	statcod := c.Ctx.ResponseWriter.Status
	statusCode := strconv.Itoa(statcod)

	// var data = recommendationData{
	// 	IdentityAddress: identityAddress,
	// 	PhoneNumber:     phoneNumber,
	// 	DomicileAddress: domicileAddress,
	// 	BusinessType:    nil,
	// }

	// fmt.Println(data)

	var response = GetRecommendationResponse{
		ResponseCode:        "0" + statusCode,
		ResponseDescription: "Success",
		Data:                &data,
	}
	c.Data["json"] = response

	c.ServeJSON()

	return response

}
