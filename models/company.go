package models

type Company struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Ruc     string `json:"ruc"`
	Name           string `json:"name"`
	Manager        string `json:"manager"`
	Owner          string `json:"owner"`
	Address        string `json:"address"`
    State      bool   `json:"state" gorm:"default:'true'"`
	
	Vehicles []Vehicle `json:"vehicles"`
	Mobiles []Mobile `json:"mobiles"`
}

func (Company) TableName() string {
	return "companies"
}
