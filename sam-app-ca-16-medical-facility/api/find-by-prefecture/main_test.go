package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	txdb.Register("txdb", "mysql", "root:password@(localhost:3307)/test_db")
	m.Run()
}

// 指定された都道府県データのみ取得する
func TestGetFacilitiesWithStatistics_Prefecture(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00001','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00001',0.44)")

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00002','test','京都府','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00002',0.44)")

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00003','test','大阪府','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00003',0.44)")

	result := getFacilitiesWithStatistics(db, "東京都", "", "")

	assert.Equal(t, 1, len(result))
}

// 指定された都市コードデータのみ取得する
func TestGetFacilitiesWithStatistics_Prefecture_CityCode(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00001','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00001',0.44)")

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00002','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131015')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00002',0.44)")

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00003','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131015')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00003',0.44)")

	result := getFacilitiesWithStatistics(db, "東京都", "131016", "")

	assert.Equal(t, 1, len(result))
}

// 指定された病院種別のみ取得する
func TestGetFacilitiesWithStatistics_Prefecture_Type(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00001','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00001',0.44)")

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00002','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'OUTPATIENT','00002',0.44)")

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00003','test','東京都','test','000-0000-0000',	35.6956,139.7924,'千代田区','131016')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'EMERGENCY','000013,0.44)")

	result := getFacilitiesWithStatistics(db, "東京都", "", "HOSPITAL")

	assert.Equal(t, 1, len(result))
}
