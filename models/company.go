package models

type Company struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Ruc     string `json:"ruc"`
	Name    string `json:"name"`
	Manager string `json:"manager"`
	Owner   string `json:"owner"`
	Address string `json:"address"`
	State   bool   `json:"state" gorm:"default:'true'"`

	CategoryID uint `json:"category_id"`

	Mobiles  []Mobile  `json:"mobiles"`
    VehicleAuthorizations []VehicleAuthorized `json:"vehicle_authorizations"`
}

func (Company) TableName() string {
	return "companies"
}
