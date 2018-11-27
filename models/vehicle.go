package models

type Vehicle struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state" gorm:"default:'true'"`
    LicensePlate string `json:"license_plate"`
    FabricationYear uint `json:"fabrication_year"`
    Brand string `json:"brand"`
    Model string `json:"model"`
    Combustible string `json:"combustible"`
    Bodywork string `json:"bodywork"`
    Axis string `json:"axis"`
	Color string `json:"color"`
    Wheel string `json:"wheel"`
    Motors uint `json:"motors"`
    Cylinders uint `json:"cylinders"`
    ChassisSeries string `json:"chassis_series"`
    Seating uint `json:"seating"`
    DryWeight string `json:"dry_weight"`
    GrossWeight string `json:"gross_weight"`
    UsefulLoad string `json:"useful_load"`
    Length float32 `json:"length"`
    Height float32 `json:"height"`
    Width float32 `json:"width"`

    Pinks                 []Pink              `json:"pinks"`
	VehicleAuthorizations []VehicleAuthorized `json:"vehicle_authorizations"`
}
