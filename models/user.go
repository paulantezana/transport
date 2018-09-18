package models

type User struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	UserName    string `json:"user_name" gorm:"type:varchar(64); unique; not null"`
	Password    string `json:"password" gorm:"type:varchar(64); not null"`
	OldPassword string `json:"old_password" gorm:"-"`
	Email       string `json:"email" gorm:"type:varchar(64); unique; not null"`
	Avatar      string `json:"avatar"`
	Profile     string `json:"profile" gorm:"type:varchar(64)"` // admin - company - municipality if company -> company_id
	CompanyID   uint `json:"company_id"`
	Key         string `json:"key"`
	State       bool   `json:"state" gorm:"default:'true'"`
}
