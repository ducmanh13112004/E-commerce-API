package models

import "time"

type ProgramDirectScreenMobileTb struct {
	RedirectId         string     `json:"redirect_id" `
	RedirectName       string     `json:"redirect_name" validate:"required"`
	TypeLink           string     `json:"type_link"`
	Actiontype         string     `json:"action_type" validate:"required" gorm:"column:action_type"`
	DataAction         string     `json:"data_action" validate:"required"`
	NameLink           string     `json:"name_link"`
	NameLinkVersionOld string     `json:"name_link_version_old"`
	Note               string     `json:"note"`
	State              int        `json:"state" validate:"oneof=0 1"`
	UpdateBy           string     `json:"update_by" `
	TUpdate            *time.Time `json:"t_update" `
}

type DeleteProgramDirectRequest struct {
	RedirectId string `json:"redirect_id" validate:"required"`
	State      int    `json:"state" validate:"oneof=0 1"`
}

func (ProgramDirectScreenMobileTb) TableName() string { return "program_direct_screen_mobile_tb1" }

func (ProgramDirectScreenMobileTb) ColumnActiontype() string         { return "action_type" }
func (ProgramDirectScreenMobileTb) ColumnDataAction() string         { return "data_action" }
func (ProgramDirectScreenMobileTb) ColumnRedirectId() string         { return "redirect_id" }
func (ProgramDirectScreenMobileTb) ColumnRedirectName() string       { return "redirect_name" }
func (ProgramDirectScreenMobileTb) ColumnTypeLink() string           { return "type_link" }
func (ProgramDirectScreenMobileTb) ColumnNameLink() string           { return "name_link" }
func (ProgramDirectScreenMobileTb) ColumnNameLinkVersionOld() string { return "name_link_version_old" }
func (ProgramDirectScreenMobileTb) ColumnNote() string               { return "note" }
func (ProgramDirectScreenMobileTb) ColumnState() string              { return "state" }
func (ProgramDirectScreenMobileTb) ColumnUpdateBy() string           { return "update_by" }
func (ProgramDirectScreenMobileTb) ColumnTUpdate() string            { return "t_update" }

type RedirectInfoProgram struct {
	RedirectId   string `json:"redirect_id"`
	RedirectName string `json:"redirect_name"`
	TypeLink     string `json:"type_link"`
	NameLink     string `json:"name_link"`
	State        int    `json:"state"`
}

type FilterProgramDirectScreenMobileTb struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}
