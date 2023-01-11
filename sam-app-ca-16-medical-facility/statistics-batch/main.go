package main

import (
	"covid19/models"
	"covid19/shared"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

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

		hospitalSubmissions := make([]models.Submission, 0)
		outpatientSubmissions := make([]models.Submission, 0)
		emergencySubmissions := make([]models.Submission, 0)

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

func getFacilities(db *sql.DB) []models.Facility {
	rows, err := db.Query("select * from Facility")
	if err != nil {
		log.Fatal(err)
	}

	var facilities []models.Facility

	for rows.Next() {
		facility := models.Facility{}
		rows.Scan(&facility.Id, &facility.Name, &facility.Prefecture, &facility.Address, &facility.Tel, &facility.Latitude, &facility.Longitude, &facility.City, &facility.CityCode)
		facilities = append(facilities, facility)
	}

	return facilities
}

func getSubmissionsByFacilityId(db *sql.DB, facilityId string) []models.Submission {
	rows, err := db.Query("select * from FacilitySubmission where facilityId = '" + facilityId + "'")
	if err != nil {
		log.Fatal(err)
	}

	var submissions []models.Submission

	for rows.Next() {
		submission := models.Submission{}
		rows.Scan(&submission.Id, &submission.Date, &submission.Answer, &submission.FacilityType, &submission.FacilityId)
		submissions = append(submissions, submission)
	}

	return submissions
}

func upsertStatistics(db *sql.DB, validDays int, normalDays int, limittedDays int, stoppedDays int, facilityType string, facilityId string, rate float64) {
	insert, err := db.Prepare(
		`insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (?,?,?,?,?,?,?) 
			on duplicate key update validDays = '` + strconv.Itoa(validDays) +
			"', normalDays = '" + strconv.Itoa(normalDays) +
			"', limittedDays = '" + strconv.Itoa(limittedDays) +
			"', stoppedDays = '" + strconv.Itoa(stoppedDays) +
			"', rate = '" + strconv.FormatFloat(rate, 'f', 2, 64) +
			"'")

	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()

	_, err = insert.Exec(validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate)
	if err != nil {
		log.Fatal(err)
	}
}

func aggregateSubmissions(submissions []models.Submission) (validDays int, normalDays int, limittedDays int, stoppedDays int, rate float64) {
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
		raw := decimal.NewFromFloat32((float32(normal) + 0.3*float32(limitted)) / float32(valid))
		rate, _ = raw.Truncate(2).Float64()
	}

	return valid, normal, limitted, stopped, rate
}
