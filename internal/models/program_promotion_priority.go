package models

import "time"

type ProgramPromotionPriority struct {
	Id               int        `gorm:"primaryKey;column:id" json:"id"`
	Priority         int        `gorm:"column:priority" json:"priority"`
	PriorityCategory string     `gorm:"column:priority_category" json:"priority_category"`
	ActiveTime       int        `gorm:"column:active_time" json:"active_time"`
	StartTime        *time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime          *time.Time `gorm:"column:end_time" json:"end_time"`
	State            int        `gorm:"column:state" json:"state"`
	MetaData         *string    `gorm:"column:meta_data" json:"meta_data"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (ProgramPromotionPriority) TableName() string              { return "program_promotion_priority" }
func (ProgramPromotionPriority) ColumnId() string               { return "id" }
func (ProgramPromotionPriority) ColumnPriority() string         { return "priority" }
func (ProgramPromotionPriority) ColumnPriorityCategory() string { return "priority_category" }
func (ProgramPromotionPriority) ColumnActiveTime() string       { return "active_time" }
func (ProgramPromotionPriority) ColumnStartTime() string        { return "start_time" }
func (ProgramPromotionPriority) ColumnEndTime() string          { return "end_time" }
func (ProgramPromotionPriority) ColumnState() string            { return "state" }
func (ProgramPromotionPriority) ColumnMetaData() string         { return "meta_data" }
func (ProgramPromotionPriority) ColumnCreatedAt() string        { return "created_at" }
func (ProgramPromotionPriority) ColumnUpdatedAt() string        { return "updated_at" }

type PriorityResponse struct {
	Id         int     `json:"id"`
	Priority   int     `json:"priority" validate:"required"`
	ActiveTime int     `json:"active_time"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	State      int     `json:"state"`
	MetaData   *string `json:"meta_data"`
}
