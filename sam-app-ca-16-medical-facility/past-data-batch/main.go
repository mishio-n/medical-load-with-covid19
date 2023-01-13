package main

import (
	"fmt"
	"log"
	"time"

	"covid19/models"
	"covid19/shared"
)

func main() {
	date := time.Date(2021, 1, 27, 0, 0, 0, 0, time.Local)
	today := time.Now().Format("20060102")

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	for {
		dateStr := date.Format("20060102")
		fmt.Println(dateStr)
		if dateStr == today {
			fmt.Println("最新までデータ取得したため終了します")
			return
		}

		dataUrl := getDataUrl(dateStr)

		dailySurveyResponse := models.FetchDailySurveyApi(dataUrl)
		models.InsertDailyServeyData(db, dailySurveyResponse)

		date = date.AddDate(0, 0, 1)
	}
}

func getDataUrl(date string) string {
	return "https://opendata.corona.go.jp/api/Covid19DailySurvey/" + date + "?prefName=%E6%9D%B1%E4%BA%AC%E9%83%BD"
}
