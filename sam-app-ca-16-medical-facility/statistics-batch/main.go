package main

import (
	"covid19/shared"
	"database/sql"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type Facility struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Prefecture string `json:"prefecture"`
	Address    string `json:"address"`
	Tel        string `json:"tel"`
	Latitude   string `json:"latitude"`
	Longtitude string `json:"longitude"`
	City       string `json:"city"`
	CityCode   string `json:"cityCode"`
}

type Submission struct {
	Id           string `json:"id"`
	Date         string `json:"date"`
	Answer       string `json:"answer"`
	FacilityType string `json:"facilityType"`
	FacilityId   string `json:"facilityId"`
}

func main() {
	db, err := shared.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute)

	facilities := getFacilities(db)
	for _, facility := range facilities {
		submissions := getSubmissionsByFacilityId(db, facility.Id)

		hospitalSubmissions := make([]Submission, 0)
		outpatientSubmissions := make([]Submission, 0)
		emergencySubmissions := make([]Submission, 0)

		for _, submission := range submissions {
			if submission.FacilityType == "HOSPITAL" {
				hospitalSubmissions = append(hospitalSubmissions, submission)
			}
			if submission.FacilityType == "OUTPATIENT" {
				outpatientSubmissions = append(outpatientSubmissions, submission)
			}
			if submission.FacilityType == "EMERGENCY" {
				emergencySubmissions = append(emergencySubmissions, submission)
			}
		}

		validDays, normalDays, limittedDays, stoppedDays, rate := aggregateSubmissions(hospitalSubmissions)
		upsertStatistics(db, validDays, normalDays, limittedDays, stoppedDays, "HOSPITAL", facility.Id, rate)

		validDays, normalDays, limittedDays, stoppedDays, rate = aggregateSubmissions(outpatientSubmissions)
		upsertStatistics(db, validDays, normalDays, limittedDays, stoppedDays, "OUTPATIENT", facility.Id, rate)

		validDays, normalDays, limittedDays, stoppedDays, rate = aggregateSubmissions(emergencySubmissions)
		upsertStatistics(db, validDays, normalDays, limittedDays, stoppedDays, "EMERGENCY", facility.Id, rate)
	}
}

func getFacilities(db *sql.DB) []Facility {
	rows, err := db.Query("select * from Facility")
	if err != nil {
		log.Fatal(err)
	}

	var facilities []Facility

	for rows.Next() {
		facility := Facility{}
		rows.Scan(&facility.Id, &facility.Name, &facility.Prefecture, &facility.Address, &facility.Tel, &facility.Latitude, &facility.Longtitude, &facility.City, &facility.CityCode)
		facilities = append(facilities, facility)
	}

	return facilities
}

func getSubmissionsByFacilityId(db *sql.DB, facilityId string) []Submission {
	rows, err := db.Query("select * from FacilitySubmission where facilityId = '" + facilityId + "'")
	if err != nil {
		log.Fatal(err)
	}

	var submissions []Submission

	for rows.Next() {
		submission := Submission{}
		rows.Scan(&submission.Id, &submission.Date, &submission.Answer, &submission.FacilityType, &submission.FacilityId)
		submissions = append(submissions, submission)
	}

	return submissions
}

func upsertStatistics(db *sql.DB, validDays int, normalDays int, limittedDays int, stoppedDays int, facilityType string, facilityId string, rate float64) {
	insert, err := db.Prepare(`insert into MedicalStatistics 
														(validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) 
														values (?,?,?,?,?,?,?) 
														on duplicate key update facilityId = '` + facilityId + "', facilityType = '" + facilityType + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()

	_, err = insert.Exec(validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate)
	if err != nil {
		log.Fatal(err)
	}
}

func aggregateSubmissions(submissions []Submission) (validDays int, normalDays int, limittedDays int, stoppedDays int, rate float64) {
	valid, normal, limitted, stopped := 0, 0, 0, 0
	rate = 0.0

	for _, submission := range submissions {
		switch submission.Answer {
		case "NORMAL":
			normal++
		case "LIMITTED":
			limitted++
		case "STOPPED":
			stopped++
		default:
			continue
		}

		valid++
	}

	if valid != 0 {
		raw := decimal.NewFromFloat32((float32(normal) + 0.5*float32(limitted)) / float32(valid))
		rate, _ = raw.Truncate(2).Float64()
	}

	return valid, normal, limitted, stopped, rate
}
