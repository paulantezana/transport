package models

type VehicleAuthorized struct {
    ID    uint   `json:"id" gorm:"primary_key"`
    CompanyID uint `json:"company_id"`
    VehicleID uint `json:"vehicle_id"`
    Authorized bool `json:"authorized"`
}

func (VehicleAuthorized) TableName() string {
    return "vehicle_authorizations"
}
