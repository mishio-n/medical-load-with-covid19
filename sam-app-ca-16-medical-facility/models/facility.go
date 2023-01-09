package models

type Facility struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Prefecture string  `json:"prefecture"`
	Address    string  `json:"address"`
	Tel        string  `json:"tel"`
	Latitude   float64 `json:"latitude"`
	Longtitude float64 `json:"longitude"`
	City       string  `json:"city"`
	CityCode   string  `json:"cityCode"`
}

type Submission struct {
	Id           string `json:"id"`
	Date         string `json:"date"`
	Answer       string `json:"answer"`
	FacilityType string `json:"facilityType"`
	FacilityId   string `json:"facilityId"`
}

type MedicalStatistics struct {
	ValidDays    int
	NormalDays   int
	LimittedDays int
	StoppedDays  int
	Rate         float64
	FacilityType string
	FacilityId   string
}
