package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"covid19/shared"
)

type FacilityWithStatistics struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Prefecture   string  `json:"prefecture"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longtitude   float64 `json:"longitude"`
	City         string  `json:"city"`
	CityCode     string  `json:"cityCode"`
	ValidDays    int     `json:"validDays"`
	NormalDays   int     `json:"normalDays"`
	LimittedDays int     `json:"limittedDays"`
	StoppedDays  int     `json:"stoppedDays"`
	Rate         float64 `json:"rate"`
	FacilityType string  `json:"facilityType"`
}

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

func getFacilitiesWithStatistics(db *sql.DB, prefecture string, cityCode string, facilityType string) []FacilityWithStatistics {
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

	var facilities []FacilityWithStatistics

	for rows.Next() {
		facility := FacilityWithStatistics{}
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
