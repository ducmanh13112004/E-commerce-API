package models

type ProgramProductTb struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	ProgramID int64  `gorm:"column:program_id"`
	SKU       string `gorm:"column:sku"`
	State     bool   `gorm:"column:state"`
	Discount  int    `gorm:"column:discount"`
}

func (ProgramProductTb) TableName() string {
	return "program_product_tb"
}

type InsertListSku struct {
	ProgramID int64 `json:"program_id" validate:"required"`
	ListSku   []struct {
		Sku             string `json:"sku" validate:"required"`
		DiscountVoucher int    `json:"discount_voucher" validate:"required"`
	} `json:"list_sku"`
}
