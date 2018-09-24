package models

type Category struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	Name             string    `json:"name"`
	CategoryParentID uint      `json:"category_parent_id"`
	Companies        []Company `json:"companies"`
}

func (Category) TableName() string {
	return "categories"
}
