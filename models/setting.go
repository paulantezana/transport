package models

type Setting struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Company        string `json:"company"`
	Email          string `json:"email"`
	Identification string `json:"identification"`
	Logo           string `json:"logo"`
	Address        string `json:"address"`
	Item           uint   `json:"item"`
}
