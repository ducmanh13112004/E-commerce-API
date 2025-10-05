package models

import (
	"time"
)

type ProgramCategory struct {
	Id         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ProgramId  int64      `json:"program_id" validate:"required"`
	CategoryId int        `json:"category_id" validate:"required"`
	Active     int        `gorm:"default:0" json:"active" validate:"oneof=0 1"`
	UpdateBy   string     `gorm:"type:varchar(500)" json:"update_by"`
	Priority   int        `json:"priority" validate:"required"`
	TUpdate    *time.Time `gorm:"type:datetime;autoUpdateTime" json:"t_update"`
	//  khóa ngoại
	ProgramPromotionTb *ProgramPromotionTb `gorm:"foreignKey:ProgramId;references:ProgramId" json:"program_promotion"`
	CategoryTb         *CategoryTb         `gorm:"foreignKey:CategoryId;references:CategoryId" json:"category"`
}

func (ProgramCategory) TableName() string          { return "program_category" }
func (ProgramCategory) ColumnId() string           { return "id" }
func (ProgramCategory) ColumnProgramId() string    { return "program_id" }
func (ProgramCategory) ColumnCategoryId() string   { return "category_id" }
func (ProgramCategory) ColumnPriority() string     { return "priority" }
func (ProgramCategory) ColumnActive() string       { return "active" }
func (ProgramCategory) ColumnCategoryname() string { return "category_name" }
func (ProgramCategory) ColumnUpdateBy() string     { return "update_by" }
func (ProgramCategory) ColumnTUpdate() string      { return "t_update" }

type CreateProgramCategory struct {
	ProgramId  int64 `json:"program_id" validate:"required"`
	CategoryId int   `json:"category_id" validate:"required"`
	Active     int   `json:"active" validate:"oneof=0 1"`
	Priority   int   `json:"priority" validate:"required"`
}

type FilterProgramcategory struct {
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	CategoryId *int64 `json:"category_id"`
	ProgramId  *int64 `json:"program_id,omitempty"`
	StateTime  *int   `json:"state_time,omitempty" validate:"omitempty,oneof=1 2 3 4 5"`
}

type UpdateOnOffProgramCategory struct {
	Id     int `gorm:"primaryKey;autoIncrement" json:"id"`
	Active int `gorm:"default:0" json:"active" validate:"oneof=0 1"`
}
