package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"covid19/models"
	"covid19/shared"
)

type FacilityWithRate struct {
	models.FacilityWithStatistics
	Distance float64 `json:"distance"`
}

type Response struct {
	AreaRate   float64            `json:"areaRate"`
	Facilities []FacilityWithRate `json:"facilities"`
	GeoJson    []Coordinates      `json:"geoJson"`
}

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

type AreaGeoJson struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Type          string  `json:"type"`
			StrokeColor   string  `json:"strokeColor"`
			StrokeOpacity int     `json:"strokeOpacity"`
			StrokeWeight  int     `json:"strokeWeight"`
			FillColor     string  `json:"fillColor"`
			FillOpacity   float64 `json:"fillOpacity"`
			MarkerColor   string  `json:"marker-color"`
		} `json:"properties,omitempty"`
		Geometry struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	longitude, err1 := strconv.ParseFloat(request.QueryStringParameters["lon"], 64)
	latitude, err2 := strconv.ParseFloat(request.QueryStringParameters["lat"], 64)

	if err1 != nil || err2 != nil {
		return events.APIGatewayProxyResponse{
			Body:       "緯度経度を正しく指定してください",
			StatusCode: 400,
		}, nil
	}

	distance, err := strconv.ParseFloat(request.QueryStringParameters["distance"], 64)
	if err != nil {
		// デフォルト値=山手線の平均駅間距離1.2km
		distance = 1.2
	}

	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	facilities := getFacilitiesArround(db, longitude, latitude, distance)
	if len(facilities) == 0 {
		return events.APIGatewayProxyResponse{
			Body:       "指定した範囲に存在するデータが見つかりませんでした",
			StatusCode: 404,
		}, nil
	}

	circlePoints := generateCirclePoints(Coordinates{Longitude: longitude, Latitude: latitude}, distance)

	response := new(Response)
	response.Facilities = facilities
	response.AreaRate = calcAreaRate(facilities)
	response.GeoJson = circlePoints

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getFacilitiesArround(db *sql.DB, longitude float64, latitude float64, distance float64) []FacilityWithRate {
	rows, err := db.Query(`select Facility.id, Facility.name, Facility.prefecture, Facility.address, Facility.latitude, Facility.longitude, Facility.city, Facility.cityCode, 
												MedicalStatistics.validDays, MedicalStatistics.normalDays, MedicalStatistics.limittedDays, MedicalStatistics.stoppedDays, MedicalStatistics.rate, MedicalStatistics.facilityType, 
												(
													6371 * acos(
														cos(radians(` + strconv.FormatFloat(latitude, 'f', 4, 64) + `))
														* cos(radians(latitude))
														* cos(radians(longitude) - radians(` + strconv.FormatFloat(longitude, 'f', 4, 64) + `))
														+ sin(radians(` + strconv.FormatFloat(latitude, 'f', 4, 64) + `))
														* sin(radians(latitude))
													)
												) AS distance 
												from Facility inner join MedicalStatistics on MedicalStatistics.facilityId=Facility.id
												where validDays >= 365
												having distance <= ` + strconv.FormatFloat(distance, 'f', 2, 64) + `
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
			&facility.Latitude,
			&facility.Longitude,
			&facility.City,
			&facility.CityCode,
			&facility.ValidDays,
			&facility.NormalDays,
			&facility.LimittedDays,
			&facility.StoppedDays,
			&facility.Rate,
			&facility.FacilityType,
			&facility.Distance,
		)

		// m単位まで桁を落とす
		facility.Distance = math.Floor(facility.Distance*1000) / 1000

		facilities = append(facilities, facility)
	}

	return facilities
}

func calcAreaRate(facilities []FacilityWithRate) float64 {
	sum := 0.0

	for _, facility := range facilities {
		sum += facility.Rate
	}

	return math.Floor(sum/float64(len(facilities))*100) / 100
}

// 指定エリアの円を作成するための座標群を作成する
// @see https://www.nanchatte.com/map/circle.html
func generateCirclePoints(center Coordinates, radius float64) []Coordinates {
	// 赤道半径(m) (WGS-84)
	const EQUATORIAL_RADIUS = 6378137

	// 扁平率の逆数 : 1/f (WGS-84)
	const F = 298.257223

	// 離心率の２乗
	E := ((2 * F) - 1) / math.Pow(F, 2)

	// 赤道半径 × π
	const PI_ER = math.Pi * EQUATORIAL_RADIUS

	// 1 - e^2 sin^2 (θ)
	TMP := 1 - E*math.Pow(math.Sin(center.Latitude*math.Pi/180), 2)

	// 経度１度あたりの長さ(m)
	arc_lat := (PI_ER * (1 - E)) / (180 * math.Pow(TMP, 3/2))

	// 緯度１度あたりの長さ(m)
	arc_lon := (PI_ER * math.Cos(center.Latitude*math.Pi/180)) / (180 * math.Pow(TMP, 1/float64(2)))

	// 半径をｍ単位に
	R := radius * 1000

	var points []Coordinates
	for i := 0; i <= 360; i++ {
		rad := float64(i) / 180 * math.Pi
		lat := (R/arc_lat)*math.Sin(rad) + center.Latitude
		lon := (R/arc_lon)*math.Cos(rad) + center.Longitude
		points = append(points, Coordinates{Latitude: lat, Longitude: lon})
	}

	return points
}

func formatGeoJson(circlePoints []Coordinates, facilityPoints []Coordinates) {

}
