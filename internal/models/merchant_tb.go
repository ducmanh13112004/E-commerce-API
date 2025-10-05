package models

import "time"

type MerchantTb struct {
	MerchantId      string     `gorm:"primaryKey" json:"merchant_id"`
	CategoryId      int        `json:"category_id"`
	MerchantName    *string    `json:"merchant_name"`
	MerchantSrcImg  *string    `json:"merchant_src_img"`
	ContactName     *string    `json:"contact_name"`
	ContactPhone    *string    `json:"contact_phone"`
	ContactPosition *string    `json:"contact_position"`
	OfficeAddress   *string    `json:"office_address"`
	Activate        int        `json:"activate"`
	TCreate         *time.Time `json:"t_create"`
	UpdateBy        *string    `json:"update_by"`
	Info            *string    `json:"info"`
}

func (MerchantTb) TableName() string { return "merchant_tb" }

func (MerchantTb) ColumnMerchantId() string      { return "merchant_id" }
func (MerchantTb) ColumnCategoryId() string      { return "category_id" }
func (MerchantTb) ColumnMerchantName() string    { return "merchant_name" }
func (MerchantTb) ColumnMerchantSrcImg() string  { return "merchant_src_img" }
func (MerchantTb) ColumnContactName() string     { return "contact_name" }
func (MerchantTb) ColumnContactPhone() string    { return "contact_phone" }
func (MerchantTb) ColumnContactPosition() string { return "contact_position" }
func (MerchantTb) ColumnOfficeAddress() string   { return "office_address" }
func (MerchantTb) ColumnActivate() string        { return "activate" }
func (MerchantTb) ColumnTCreate() string         { return "t_create" }
func (MerchantTb) ColumnUpdateBy() string        { return "update_by" }
func (MerchantTb) ColumnInfo() string            { return "info" }
