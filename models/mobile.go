package models

type Mobile struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	MacAddress string `json:"mac_address"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	Driver     string `json:"driver"`
	State      bool   `json:"state" gorm:"default:'true'"`
}
