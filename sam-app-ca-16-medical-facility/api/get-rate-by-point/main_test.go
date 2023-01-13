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

// 指定エリアに含まれるレコードがあった場合にデータを返す
func TestGetFacilitiesArround_Success(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00001','test','test','test','000-0000-0000',	35.6956,139.7924,'test','test')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00001',0.44)")

	result := getFacilitiesArround(db, 139.7711, 35.6916, 2)

	assert.Equal(t, 1, len(result))
}

// 指定エリアに含まれるが、データ数が1年未満の場合はレコードを返さない
func TestGetFacilitiesArround_LackDays(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00001','test','test','test','000-0000-0000',	35.6956,139.7924,'test','test')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (364,100,200,60,'HOSPITAL','00001',0.44)")

	result := getFacilitiesArround(db, 139.7711, 35.6916, 2)

	assert.Equal(t, 0, len(result))
}

// 指定エリアに含まれるレコードがない場合は空データを返す
func TestGetFacilitiesArround_NotFound(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("insert into Facility (id,name,prefecture,address,tel,latitude,longitude,city,cityCode) VALUES ('00001','test','test','test','000-0000-0000',	35.6956,139.7924,'test','test')")
	db.Exec("insert into MedicalStatistics (validDays, normalDays, limittedDays, stoppedDays, facilityType, facilityId, rate) values (365,100,200,60,'HOSPITAL','00001',0.44)")

	result := getFacilitiesArround(db, 139.7711, 35.6916, 1)

	assert.Equal(t, 0, len(result))
}
