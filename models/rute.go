package models

type Rute struct {
    ID    uint `json:"id" gorm:"primary_key"`
    State bool `json:"state" gorm:"default:'true'"`
}
