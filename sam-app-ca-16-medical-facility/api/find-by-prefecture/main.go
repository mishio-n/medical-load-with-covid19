package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"covid19/models"
	"covid19/shared"
)

type Request struct {
	Prefecture string `json:"prefecture"`
	CityCode   string `json:"cityCode"`
}

type Response struct {
	models.Facility
	models.MedicalStatistics
	Id string `json:"facilityId"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prefecture := request.QueryStringParameters["prefecture"]
	cityCode := request.QueryStringParameters["cityCode"]

	if prefecture == "" {
		return events.APIGatewayProxyResponse{
			Body:       "県名を指定してください",
			StatusCode: 400,
		}, nil
	}

	var response Response
	log.Println(response)

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	facilities := getFacilities(db, prefecture, cityCode)
	body, _ := json.Marshal(facilities)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getFacilities(db *sql.DB, prefecture string, cityCode string) []models.Facility {
	statement := "select * from Facility where prefecture = '" + prefecture + "'"
	if cityCode != "" {
		statement += " and cityCode = '" + cityCode + "'"
	}

	rows, err := db.Query(statement)
	if err != nil {
		log.Fatal(err)
	}

	var facilities []models.Facility

	for rows.Next() {
		facility := models.Facility{}
		rows.Scan(&facility.Id, &facility.Name, &facility.Prefecture, &facility.Address, &facility.Tel, &facility.Latitude, &facility.Longtitude, &facility.City, &facility.CityCode)
		facilities = append(facilities, facility)
	}

	return facilities
}
