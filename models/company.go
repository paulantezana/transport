package models

type Company struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	Manager        string `json:"manager"`
	Owner          string `json:"owner"`
	Identification string `json:"identification"`
	Address        string `json:"address"`
    State      bool   `json:"state" gorm:"default:'true'"`
}

func (Company) TableName() string {
	return "companies"
}
