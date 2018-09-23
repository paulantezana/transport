package models

import "time"

type Pink struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	MobileID  uint      `json:"mobile_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	DatePiker time.Time `json:"date_piker"`
	
	VehicleID uint `json:"vehicle_id"`
}
