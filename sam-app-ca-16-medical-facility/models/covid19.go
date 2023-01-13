package models

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
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

func FetchDailySurveyApi(url string) []DailySurveyResponse {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
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

	return dailySurveyResponse
}

func InsertDailyServeyData(db *sql.DB, dailySurveyResponse []DailySurveyResponse) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	facilityInsert, err := tx.Prepare("insert ignore into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES (?,?,?,?,?,?,?,?,?)")
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
		longitude, err := strconv.ParseFloat(row.Longitude, 64)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := facilityInsert.Exec(row.FacilityId, row.FacilityName, row.PrefName, row.FacilityAddr, row.FacilityTel, latitude, longitude, row.CityName, row.LocalGovCode); err != nil {
			log.Fatal(err)
		}

		if _, err := submissionInsert.Exec(row.SubmitDate, convertAnsType(row.AnsType), convertFacilityType(row.FacilityType), row.FacilityId); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func convertFacilityType(raw string) string {
	switch raw {
	case "入院":
		return "HOSPITAL"
	case "外来":
		return "OUTPATIENT"
	case "救急":
		return "EMERGENCY"
	default:
		return ""
	}
}

func convertAnsType(raw string) string {
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
