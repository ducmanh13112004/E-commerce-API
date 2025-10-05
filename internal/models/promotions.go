package models

import "time"

type PromotionTb struct {
	PromotionCode      string     `gorm:"primaryKey" json:"promotion_code"`
	ProgramId          int64      `gorm:"primaryKey" json:"program_id"`
	TotalRelease       int        `json:"total_release"`
	TotalCurrent       int        `json:"total_current"`
	TotalUse           int        `json:"total_use"`
	StatePromotion     int        `json:"state_promotion"`
	TCreate            *time.Time `json:"t_create"`
	PromotionEndUsable *time.Time `json:"promotion_end_usable"`
}

func (PromotionTb) TableName() string { return "promotion_tb" }

func (PromotionTb) ColumnPromotionCode() string      { return "promotion_code" }
func (PromotionTb) ColumnProgramId() string          { return "program_id" }
func (PromotionTb) ColumnTotalRelease() string       { return "total_release" }
func (PromotionTb) ColumnTotalCurrent() string       { return "total_current" }
func (PromotionTb) ColumnTotalUse() string           { return "total_use" }
func (PromotionTb) ColumnStatePromotion() string     { return "state_promotion" }
func (PromotionTb) ColumnTCreate() string            { return "t_create" }
func (PromotionTb) ColumnPromotionEndUsable() string { return "promotion_end_usable" }

type InputReceiveVoucher struct {
	CustomerId    string `json:"customer_id" validate:"required"`
	CustomerPhone string `json:"customer_phone" validate:"required"`
	ProgramId     int64  `json:"program_id" validate:"required"`
}

type RespSPReceiveVoucher struct {
	ProgramInfo          *ProgramPromotionTb
	VoucherInfo          *PromotionTb
	AvailableRemainTimes int
}

type InputSPReceiceVoucher struct {
	ProgramId      int64
	CustomerId     string
	CustomerPhone  string
	ProgramTb      *ProgramPromotionTb
	ReceiveCurrent int // số lượng user sẽ nhận được
	StatePromotion int // trạng thái voucher
	StateUsable    int // trạng thái sử dụng voucher
}

type InputAfiliateGetProgramListFrCategoryId struct {
	CategoryId int `json:"category_id" validate:"required"`
}

type TotalRemainPromotionProgramId struct {
	TotalCurrent int `json:"total_current"`
	TotalUse     int `json:"total_use"`
	TotalRelease int `json:"total_release"`
}

type CreatePromotion struct {
	PromotionCode      string `json:"promotion_code" validate:"required"`
	ProgramId          int64  `json:"program_id" validate:"required"`
	TotalRelease       int    `json:"total_release" validate:"required"`
	PromotionEndUsable string `json:"promotion_end_usable" validate:"required"`
}

type InputCreatePromotion struct {
	ProgramId int64 `json:"program_id" validate:"required"`
	Type      int   `json:"type" validate:"required,oneof=0 1 2"`
	Data      struct {
		Prefix             *string `json:"prefix,omitempty" validate:"omitempty,required"`
		Quantity           *int    `json:"quantity,omitempty" validate:"omitempty,required"`
		TotalRelease       *int    `json:"total_release,omitempty" validate:"omitempty,required"`
		PromotionEndUsable *string `json:"promotion_end_usable,omitempty" validate:"omitempty,required"`

		PromotionCode *string `json:"promotion_code,omitempty" validate:"omitempty,required"`
		ListPromotion []struct {
			PromotionCode      string `json:"promotion_code" validate:"required"`
			PromotionEndUsable string `json:"promotion_end_usable" validate:"required"`
		} `json:"list_promotion,omitempty"`
	} `json:"data"`
}
