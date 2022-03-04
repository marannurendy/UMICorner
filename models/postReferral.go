package models

import (
	_ "beego/conf"
	_ "fmt"
)

type TransactionUpdate struct {
	Status                    string `json:"status"`
	ExecutorId                string `json:"executorId"`
	ExecutorName              string `json:"executorName"`
	ExecutorPhoneNumber       string `json:"executorPhoneNumber"`
	ExecutionBusinessUnitId   string `json:"executionBusinessUnitId"`
	ExecutionBusinessUnitName string `json:"executionBusinessUnitName"`
	TransactionId             string `json:"transactionId"`
}

type TransactionUpdateResponse struct {
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
}

type TransactionPost struct {
	Nik              string `json:"nik"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phoneNumber"`
	Birthdate        string `json:"birthDate"`
	Gender           string `json:"gender"`
	Education        string `json:"education"`
	Address          string `json:"address"`
	RtRw             string `json:"rtRw"`
	Village          string `json:"village"`
	District         string `json:"district"`
	Province         string `json:"province"`
	PostalCode       string `json:"postalCode"`
	Profession       string `json:"profession"`
	BusinessType     string `json:"businessType"`
	Branch           string `json:"branch"`
	LoanLimit        string `json:"loanLimit"`
	LoanPeriod       string `json:"loanPeriod"`
	PartnerProductId string `json:"partnerProductId"`
	IdempotencyKey   string `json:"idempotencyKey"`
}

type TransactionResponse struct {
	ResponseCode        string  `json:"responseCode"`
	ResponseDescription string  `json:"responseDescription"`
	Data                DetData `json:"data"`
}

type DetData struct {
	Transactionid string `json:"transactionId"`
}

type InquiryPost struct {
	TransactionId string `json:"TransactionId"`
}

type TransactionIncResponse struct {
	ResponseCode        string       `json:"responseCode"`
	ResponseDescription string       `json:"responseDescription"`
	Data                *DataInquiry `json:"data"`
}

type DataInquiry struct {
	Status                    string `json:"status"`
	TransactionId             string `json:"transactionId"`
	ExecutorId                string `json:"executorId"`
	ExecutorName              string `json:"executorName:"`
	ExecutorPhoneNumber       string `json:"executorPhoneNumber"`
	ExecutionBusinessUnitId   string `json:"executionBusinessUnitId"`
	ExecutionBusinessUnitName string `json:"executionBusinessUnitName"`
	LoanRealization           string `json:"loanRealization"`
}

var Apikey = "QlJJUE5NSmF5YVNlbGFsdQ=="
