package models

type DailySurveyResponse struct {
	FacilityId   string `json:"facilityId"`
	FacilityName string `json:"facilityName"`
	ZipCode      string `json:"zipCode"`
	PrefName     string `json:"prefName"`
	FacilityAddr string `json:"facilityAddr"`
	FacilityTel  string `json:"facilityTel"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	SubmitDate   string `json:"submitDate"`
	FacilityType string `json:"facilityType"`
	AnsType      string `json:"ansType"`
	LocalGovCode string `json:"localGovCode"`
	CityName     string `json:"cityName"`
	FacilityCode string `json:"facilityCode"`
}

func ConvertFacilityType(raw string) string {
	switch raw {
	case "入院":
		return "HOSPITAL"
	case "外来":
		return "OUTPATIENT"
	case "救急":
		return "EMERGENCY"
	default:
		return ""
	}
}

func ConvertAnsType(raw string) string {
	switch raw {
	case "通常":
		return "NORMAL"
	case "制限":
		return "LIMITTED"
	case "停止":
		return "STOPPED"
	case "未回答":
		return "NOANSWER"
	default:
		return "NULL"
	}
}
