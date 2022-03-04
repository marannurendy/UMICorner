package models

import (
	_ "beego/conf"
	_ "fmt"
)

type TransactionGet struct {
	TransactionId string `json:"transactionId"`
}

type GetResponse struct {
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
	Data                DataGet
}

type DataGet struct {
	Status                    string `json:"Status"`
	TransactionId             string `json:"TransactionId"`
	ExecutorId                string `json:"ExecutorId"`
	ExecutorName              string `json:"ExecutorName"`
	ExecutorPhoneNumber       string `json:"ExecutorPhoneNumber"`
	ExecutionBusinessUnitId   string `json:"ExecutionBusinessUnitId"`
	ExecutionBusinessUnitName string `json:"ExecutionBusinessUnitName"`
	LoanRealization           string `json:"LoanRealization"`
}
