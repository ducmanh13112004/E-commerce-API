package models

import "time"

type ProgramRedirectTb struct {
	RedirectId   int        `json:"redirect_id" gorm:"primaryKey;column:redirect_id"`
	RedirectName string     `json:"redirect_name" validate:"required" gorm:"column:redirect_name"`
	State        int        `json:"state" validate:"oneof=0 1" gorm:"column:state"`
	UpdateBy     string     `json:"update_by" gorm:"column:update_by"`
	TCreate      *time.Time `json:"t_create" gorm:"column:t_create"`
	TUpdate      *time.Time `json:"t_update" gorm:"column:t_update"`
}

type DeleteProgramRedirectRequest struct {
	RedirectId int `json:"redirect_id" validate:"required"`
	State      int `json:"state" validate:"oneof=0 1"`
}
type RedirectInfo struct {
	RedirectId   int    `json:"redirect_id"`
	RedirectName string `json:"redirect_name"`
}

func (ProgramRedirectTb) TableName() string {
	return "program_redirect_tb"
}

func (ProgramRedirectTb) ColumnRedirectId() string   { return "redirect_id" }
func (ProgramRedirectTb) ColumnRedirectName() string { return "redirect_name" }
func (ProgramRedirectTb) ColumnState() string        { return "state" }
func (ProgramRedirectTb) ColumnUpdateBy() string     { return "update_by" }
func (ProgramRedirectTb) ColumnTCreate() string      { return "t_create" }
func (ProgramRedirectTb) ColumnTUpdate() string      { return "t_update" }

type InfoRedirectProgram struct {
	RedirectId   int    `json:"redirect_id"`
	RedirectName string `json:"redirect_name"`
	State        int    `json:"state"`
}

// Struct dùng để filter dữ liệu
type FilterProgramRedirectTb struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}
