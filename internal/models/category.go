package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type CategoryTb struct {
	CategoryId       int        `gorm:"primaryKey" json:"category_id"`
	CategoryName     string     `json:"category_name" validate:"required,omitempty"`
	CategoryPriority int        `json:"category_priority" validate:"required,min=0,max=999999"`
	CategorySrcImg   string     `json:"category_src_img" validate:"required"`
	State            *int       `json:"state" validate:"required,oneof=0 1"`
	UpdateBy         string     `gorm:"type:varchar(500)" json:"update_by"`
	TUpdate          *time.Time `gorm:"type:datetime" json:"t_update"`

	ListProgram []ProgramCategory `gorm:"foreignKey:CategoryId;references:CategoryId" json:"program_category"`
}

type DeleteCategoryRequest struct {
	CategoryId int `json:"category_id" validate:"required"`
	State      int `json:"state" validate:"oneof=0 1"`
}

type FilterCategory struct {
	CategoryId int  `json:"category_id" `
	State      *int `json:"state" `
}

// Hàm Validate để kiểm tra dữ liệu đầu vào
func (c *CategoryTb) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (CategoryTb) TableName() string { return "category_tb" }

func (CategoryTb) ColumnCategoryId() string       { return "category_id" }
func (CategoryTb) ColumnCategoryName() string     { return "category_name" }
func (CategoryTb) ColumnCategoryPriority() string { return "category_priority" }
func (CategoryTb) ColumnCategorySrcImg() string   { return "category_src_img" }
func (CategoryTb) ColumnState() string            { return "state" }
func (CategoryTb) ColumnUpdateBy() string         { return "update_by" }
func (CategoryTb) ColumnTUpdate() string          { return "t_update" }
