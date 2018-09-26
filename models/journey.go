package models

type Journey struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	Origin          string `json:"origin"`
	Destination     string `json:"destination"`
	Distance        string `json:"distance"`
	Frequency       uint   `json:"frequency"`
	StartPoint      uint   `json:"start_point"`
	CommercialSpeed uint   `json:"commercial_speed"`
	State           bool   `json:"state" gorm:"default:'true'"`

	CompanyID uint `json:"company_id"`
	JourneyDetails []JourneyDetail `json:"journey_details"`
}
