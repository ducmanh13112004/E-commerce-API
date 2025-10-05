package models

import "time"

type HistoriesUsableTb struct {
	Id            int        `gorm:"primaryKey" json:"id"`
	CustomerID    *string    `json:"customer_id"`
	PromotionCode string     `json:"promotion_code"`
	StateUsable   int        `json:"state_usable"`
	TCreate       *time.Time `json:"t_create"`
	OrderID       *string    `json:"order_id"`
	MethodUsable  *string    `json:"method_usable"`
	UpdateBy      *string    `json:"update_by"`
	TUpdate       *time.Time `json:"t_update"`
}

func (HistoriesUsableTb) TableName() string { return "histories_usable_tb" }

func (HistoriesUsableTb) ColumnId() string            { return "id" }
func (HistoriesUsableTb) ColumnCustomerId() string    { return "customer_id" }
func (HistoriesUsableTb) ColumnPromotionCode() string { return "promotion_code" }
func (HistoriesUsableTb) ColumnStateUsable() string   { return "state_usable" }
func (HistoriesUsableTb) ColumnTCreate() string       { return "t_create" }
func (HistoriesUsableTb) ColumnOrderId() string       { return "order_id" }
func (HistoriesUsableTb) ColumnMethodUsable() string  { return "method_usable" }
func (HistoriesUsableTb) ColumnUpdateBy() string      { return "update_by" }
func (HistoriesUsableTb) ColumnTUpdate() string       { return "t_update" }

type InfoPromotionCodeCus struct {
	// này của table  histories_usable_tb
	Id            int        `gorm:"primaryKey" json:"id"`
	CustomerID    *string    `json:"customer_id"`
	PromotionCode string     `json:"promotion_code"`
	StateUsable   int        `json:"state_usable"`
	TCreate       *time.Time `json:"t_create"`
	OrderID       *string    `json:"order_id"`
	MethodUsable  *string    `json:"method_usable"`
	UpdateBy      *string    `json:"update_by"`
	TUpdate       *time.Time `json:"t_update"`
	// này của table  promotion_tb
	PromotionEndUsable *time.Time `json:"promotion_end_usable"`
}
