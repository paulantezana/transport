package models

type Vehicle struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state" gorm:"default:'true'"`

	CompanyID uint `json:"company_id"`

	Pinks []Pink `json:"pinks"`
	VehicleAuthorizations []VehicleAuthorized `json:"vehicle_authorizations"`
}
