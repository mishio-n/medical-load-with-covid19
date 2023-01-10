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
	ValidDays    int     `json:"validDays"`
	NormalDays   int     `json:"normalDays"`
	LimittedDays int     `json:"limittedDays"`
	StoppedDays  int     `json:"stoppedDays"`
	Rate         float64 `json:"rate"`
	FacilityType string  `json:"facilityType"`
	FacilityId   string  `json:"facilityId"`
}

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
