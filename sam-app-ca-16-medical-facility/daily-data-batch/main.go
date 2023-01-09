package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"covid19/shared"
)

var (
	DATA_URL = "https://opendata.corona.go.jp/api/Covid19DailySurvey?localGovCode=261009"
)

type DailySurveyResponse struct {
	FacilityId   string `json:"facilityId"`
	FacilityName string `json:"facilityName"`
	ZipCode      string `json:"zipCode"`
	PrefName     string `json:"prefName"`
	FacilityAddr string `json:"facilityAddr"`
	FacilityTel  string `json:"facilityTel"`
	Latitude     string `json:"latitude"`
	Longtitude   string `json:"longitude"`
	SubmitDate   string `json:"submitDate"`
	FacilityType string `json:"facilityType"`
	AnsType      string `json:"ansType"`
	LocalGovCode string `json:"localGovCode"`
	CityName     string `json:"cityName"`
	FacilityCode string `json:"facilityCode"`
}

func handler(request events.CloudWatchEvent) {
	req, _ := http.NewRequest(http.MethodGet, DATA_URL, nil)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error Request:", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal("Error Response:", res.Status)
	}

	body, _ := io.ReadAll(res.Body)

	var dailySurveyResponse []DailySurveyResponse
	json.Unmarshal(body, &dailySurveyResponse)

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	facilityInsert, err := tx.Prepare("insert ignore into Facility (id,name,prefecture,address,tel,latitude,longtitude,city,cityCode) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer facilityInsert.Close()

	submissionInsert, err := tx.Prepare("insert ignore into FacilitySubmission (date,answer,facilityType,facilityId) VALUES (?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer submissionInsert.Close()

	for _, row := range dailySurveyResponse {
		latitude, err := strconv.ParseFloat(row.Latitude, 64)
		if err != nil {
			log.Fatal(err)
		}
		longtitude, err := strconv.ParseFloat(row.Longtitude, 64)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := facilityInsert.Exec(row.FacilityId, row.FacilityName, row.PrefName, row.FacilityAddr, row.FacilityTel, latitude, longtitude, row.CityName, row.LocalGovCode); err != nil {
			log.Fatal(err)
		}

		if _, err := submissionInsert.Exec(row.SubmitDate, convAnsType(row.AnsType), convFacilityType(row.FacilityType), row.FacilityId); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(handler)
}

func convFacilityType(raw string) string {
	switch raw {
	case "入院":
		return "HOSPITAL"
	case "救急":
		return "OUTPATIENT"
	case "外来":
		return "EMERGENCY"
	default:
		return ""
	}
}

func convAnsType(raw string) string {
	switch raw {
	case "通常":
		return "NORMAL"
	case "制限":
		return "LIMITTED"
	case "停止":
		return "STOPPED"
	case "未回答":
		return "NOANSWER"
	default:
		return "NULL"
	}
}
