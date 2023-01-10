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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prefecture := request.QueryStringParameters["prefecture"]
	cityCode := request.QueryStringParameters["cityCode"]
	facilityType := request.QueryStringParameters["type"]

	if prefecture == "" {
		return events.APIGatewayProxyResponse{
			Body:       "県名を指定してください",
			StatusCode: 400,
		}, nil
	}

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	facilities := getFacilitiesWithStatistics(db, prefecture, cityCode, facilityType)
	body, _ := json.Marshal(facilities)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getFacilitiesWithStatistics(db *sql.DB, prefecture string, cityCode string, facilityType string) []models.FacilityWithStatistics {
	statement := `select Facility.id, Facility.name, Facility.prefecture, Facility.address, Facility.latitude, Facility.longtitude, Facility.city, Facility.cityCode, 
								MedicalStatistics.validDays, MedicalStatistics.normalDays, MedicalStatistics.limittedDays, MedicalStatistics.stoppedDays, MedicalStatistics.rate, MedicalStatistics.facilityType
								from Facility inner join MedicalStatistics on MedicalStatistics.facilityId=Facility.id 
								where prefecture = '` + prefecture + "'"
	if cityCode != "" {
		statement += " and cityCode = '" + cityCode + "'"
	}

	if facilityType != "" {
		statement += " and facilityType = '" + facilityType + "'"
	}

	rows, err := db.Query(statement)
	if err != nil {
		log.Fatal(err)
	}

	var facilities []models.FacilityWithStatistics

	for rows.Next() {
		facility := models.FacilityWithStatistics{}
		rows.Scan(
			&facility.Id,
			&facility.Name,
			&facility.Prefecture,
			&facility.Address,
			&facility.Latitude,
			&facility.Longtitude,
			&facility.City,
			&facility.CityCode,
			&facility.ValidDays,
			&facility.NormalDays,
			&facility.LimittedDays,
			&facility.StoppedDays,
			&facility.Rate,
			&facility.FacilityType,
		)

		facilities = append(facilities, facility)
	}
	return facilities
}
