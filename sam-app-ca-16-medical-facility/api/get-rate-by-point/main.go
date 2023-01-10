package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"covid19/models"
	"covid19/shared"
)

type FacilityWithRate struct {
	models.Facility
	Distance float64 `json:"distance"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	longtitude, err1 := strconv.ParseFloat(request.QueryStringParameters["lon"], 64)
	latitude, err2 := strconv.ParseFloat(request.QueryStringParameters["lat"], 64)

	if err1 != nil || err2 != nil {
		return events.APIGatewayProxyResponse{
			Body:       "緯度経度を正しく指定してください",
			StatusCode: 400,
		}, nil
	}

	distance, err := strconv.Atoi(request.QueryStringParameters["distance"])
	if err != nil {
		// デフォルト値=山手線の平均駅間距離1.2km
		distance = 1200
	}

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	facilities := getFacilitiesArround(db, longtitude, latitude, distance)
	body, _ := json.Marshal(facilities)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getFacilitiesArround(db *sql.DB, longtitude float64, latitude float64, distance int) []FacilityWithRate {
	rows, err := db.Query(`select *, 
												(
													6371 * acos(
														cos(radians(35.6916))
														* cos(radians(latitude))
														* cos(radians(longtitude) - radians(139.7703))
														+ sin(radians(35.6916))
														* sin(radians(latitude))
													)
												) AS distance 
												from Facility 
												having distance <= 1.2
												order by distance`)
	if err != nil {
		log.Fatal(err)
	}

	var facilities []FacilityWithRate
	for rows.Next() {
		facility := FacilityWithRate{}
		rows.Scan(
			&facility.Id,
			&facility.Name,
			&facility.Prefecture,
			&facility.Address,
			&facility.Tel,
			&facility.Latitude,
			&facility.Longtitude,
			&facility.City,
			&facility.CityCode,
			&facility.Distance,
		)

		facilities = append(facilities, facility)
	}

	return facilities
}

// func getFacilitiesWithStatistics(db *sql.DB, prefecture string, cityCode string, facilityType string) []models.FacilityWithStatistics {
// 	statement := `select Facility.id, Facility.name, Facility.prefecture, Facility.address, Facility.latitude, Facility.longtitude, Facility.city, Facility.cityCode,
// 								MedicalStatistics.validDays, MedicalStatistics.normalDays, MedicalStatistics.limittedDays, MedicalStatistics.stoppedDays, MedicalStatistics.rate, MedicalStatistics.facilityType
// 								from Facility inner join MedicalStatistics on MedicalStatistics.facilityId=Facility.id
// 								where prefecture = '` + prefecture + "'"
// 	if cityCode != "" {
// 		statement += " and cityCode = '" + cityCode + "'"
// 	}

// 	if facilityType != "" {
// 		statement += " and facilityType = '" + facilityType + "'"
// 	}

// 	rows, err := db.Query(statement)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var facilities []models.FacilityWithStatistics

// 	for rows.Next() {
// 		facility := models.FacilityWithStatistics{}
// 		rows.Scan(
// 			&facility.Id,
// 			&facility.Name,
// 			&facility.Prefecture,
// 			&facility.Address,
// 			&facility.Latitude,
// 			&facility.Longtitude,
// 			&facility.City,
// 			&facility.CityCode,
// 			&facility.ValidDays,
// 			&facility.NormalDays,
// 			&facility.LimittedDays,
// 			&facility.StoppedDays,
// 			&facility.Rate,
// 			&facility.FacilityType,
// 		)

// 		facilities = append(facilities, facility)
// 	}
// 	return facilities
// }
