package models

type JourneyDetail struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Name     string  `json:"name"`
	Sequence uint    `json:"sequence"`
	Latitude float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`

	JourneyID uint `json:"journey_id"`
}
