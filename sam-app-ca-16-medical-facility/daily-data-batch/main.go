package main

import (
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"covid19/models"
	"covid19/shared"
)

var (
	DATA_URL = "https://opendata.corona.go.jp/api/Covid19DailySurvey/?prefName=%E6%9D%B1%E4%BA%AC%E9%83%BD"
)

type DailySurveyResponse struct {
	FacilityId   string `json:"facilityId"`
	FacilityName string `json:"facilityName"`
	ZipCode      string `json:"zipCode"`
	PrefName     string `json:"prefName"`
	FacilityAddr string `json:"facilityAddr"`
	FacilityTel  string `json:"facilityTel"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	SubmitDate   string `json:"submitDate"`
	FacilityType string `json:"facilityType"`
	AnsType      string `json:"ansType"`
	LocalGovCode string `json:"localGovCode"`
	CityName     string `json:"cityName"`
	FacilityCode string `json:"facilityCode"`
}

func handler(request events.CloudWatchEvent) {
	dailySurveyResponse := models.FetchDailySurveyApi(DATA_URL)

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	models.InsertDailyServeyData(db, dailySurveyResponse)
}

func main() {
	lambda.Start(handler)
}
