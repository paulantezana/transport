package models

type Category struct {
	ID                 uint   `json:"id" gorm:"primary_key"`
	Name               string `json:"name"`
	CategoryParentID   uint   `json:"category_parent_id"`
	CategoryParentName string `json:"category_parent_name"`
    State      bool   `json:"state" gorm:"default:'true'"`
}

func (Category) TableName() string {
	return "categories"
}
