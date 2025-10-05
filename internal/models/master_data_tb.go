package models

import "time"

type MasterData struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code         string     `gorm:"column:code;unique;not null" json:"code"`
	DataKeyValue string     `gorm:"column:data_key_value;type:text" json:"data_key_value"`
	CreateAt     *time.Time `gorm:"column:create_at" json:"create_at,omitempty"`
	UpdateAt     *time.Time `gorm:"column:update_at" json:"update_at,omitempty"`
	UpdateBy     *string    `gorm:"column:update_by" json:"update_by,omitempty"`
	Active       int        `gorm:"column:active;default:1" json:"active"`
}

func (MasterData) TableName() string {
	return "master_data_tb"
}

type MasterDataInput struct {
	Code         string      `json:"code" validate:"required"`
	DataKeyValue interface{} `json:"data_key_value" validate:"required"`
	UpdateBy     string      `json:"update_by"`
}
