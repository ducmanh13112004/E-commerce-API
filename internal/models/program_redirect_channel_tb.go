package models

import "time"

type ProgramRedirectChannelTb struct {
	ID         int        `json:"id" gorm:"primaryKey;column:id"`
	RedirectId int        `json:"redirect_id" gorm:"column:redirect_id"`
	Channel    string     `json:"channel" gorm:"column:channel"`
	MinVersion *string    `json:"min_version" gorm:"column:min_version"` //
	MaxVersion *string    `json:"max_version" gorm:"column:max_version"`
	ActionType string     `json:"action_type" validate:"required" gorm:"column:action_type"`
	DataAction string     `json:"data_action" validate:"required" gorm:"column:data_action"`
	Data       *string    `json:"data" gorm:"column:data"`
	TUpdate    *time.Time `json:"t_update" gorm:"column:t_update"`
}

type DeleteProgramRedirectChannelRequest struct {
	RedirectId int `json:"redirect_id" validate:"required"`
}

func (ProgramRedirectChannelTb) TableName() string {
	return "program_redirect_channel_tb"
}

func (ProgramRedirectChannelTb) ColumnID() string         { return "id" }
func (ProgramRedirectChannelTb) ColumnRedirectId() string { return "redirect_id" }
func (ProgramRedirectChannelTb) ColumnChannel() string    { return "channel" }
func (ProgramRedirectChannelTb) ColumnMinVersion() string { return "min_version" }
func (ProgramRedirectChannelTb) ColumnMaxVersion() string { return "max_version" }
func (ProgramRedirectChannelTb) ColumnActionType() string { return "action_type" }
func (ProgramRedirectChannelTb) ColumnDataAction() string { return "data_action" }
func (ProgramRedirectChannelTb) ColumnData() string       { return "data" }
func (ProgramRedirectChannelTb) ColumnTUpdate() string    { return "t_update" }

type RedirectInfoProgramChannel struct {
	ID         int    `json:"id"`
	RedirectId int    `json:"redirect_id"`
	Channel    string `json:"channel"`
	MinVersion string `json:"min_version"`
	MaxVersion string `json:"max_version"`
	ActionType string `json:"action_type"`
}

type FilterProgramRedirectChannelTb struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

type ProgramRedirectJoin struct {
	RedirectId   int        `json:"redirect_id" gorm:"column:redirect_id"`
	RedirectName string     `json:"redirect_name" gorm:"column:redirect_name"`
	State        int        `json:"state" gorm:"column:state"`
	UpdateBy     string     `json:"update_by" gorm:"column:update_by"`
	TCreate      *time.Time `json:"t_create" gorm:"column:t_create"`
	TUpdate      *time.Time `json:"t_update" gorm:"column:t_update"`

	// Channel info
	Channel    string     `json:"channel" gorm:"column:channel"`
	MinVersion string     `json:"min_version" gorm:"column:min_version"`
	MaxVersion string     `json:"max_version" gorm:"column:max_version"`
	ActionType string     `json:"action_type" gorm:"column:action_type"`
	DataAction string     `json:"data_action" gorm:"column:data_action"`
	Data       string     `json:"data" gorm:"column:data"`
	CTUpdate   *time.Time `json:"channel_t_update" gorm:"column:channel_t_update"`
}

// DTO cho request body
type InputProgramRedirect struct {
	RedirectId          int                           `json:"redirect_id"`
	RedirectName        string                        `json:"redirect_name" validate:"required"`
	State               int                           `json:"state" validate:"oneof=0 1"`
	ListRedirectChannel []InputProgramRedirectChannel `json:"list_redict_channel"`
	Email               string                        `json:"-"` // mình set từ ctx.Locals
}

type InputProgramRedirectChannel struct {
	Channel    string                 `json:"channel"`
	ActionType string                 `json:"actionType"`
	Data       map[string]interface{} `json:"data"`
	DataAction string                 `json:"dataAction"`
	// Note       string                 `json:"Note"`
}
type InputGetPromotion struct {
	Data struct {
		RedirectId int `json:"redirect_id" validate:"required"`
	} `json:"data"`
}
