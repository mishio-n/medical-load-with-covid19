package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"covid19/models"
	"covid19/shared"
)

func main() {
	date := time.Date(2021, 1, 27, 0, 0, 0, 0, time.Local)
	today := time.Now().Format("20060102")

	for {
		dateStr := date.Format("20060102")
		fmt.Println(dateStr)
		if dateStr == today {
			fmt.Println("最新までデータ取得したため終了します")
			return
		}

		dataUrl := getDataUrl(dateStr)
		fetchAndInsertDailyData(dataUrl)
		date = date.AddDate(0, 0, 1)
	}
}

func fetchAndInsertDailyData(dataUrl string) {
	req, _ := http.NewRequest(http.MethodGet, dataUrl, nil)
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

	var dailySurveyResponse []models.DailySurveyResponse
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

		if _, err := submissionInsert.Exec(row.SubmitDate, convAnsType(row.AnsType), convFacilityType(row.FacilityType), row.FacilityId); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func getDataUrl(date string) string {
	return "https://opendata.corona.go.jp/api/Covid19DailySurvey/" + date + "?prefName=%E6%9D%B1%E4%BA%AC%E9%83%BD"
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
