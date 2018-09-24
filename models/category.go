package models

type Category struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	Name             string    `json:"name"`
	CategoryParentID uint      `json:"category_parent_id"`
	State            bool      `json:"state" gorm:"default:'true'"`
	Companies        []Company `json:"companies"`
}

func (Category) TableName() string {
	return "categories"
}
