package models

type ProgramLocationTb struct {
	Id         int `json:"id"`
	ProgramId  int `json:"program_id"`
	LocationId int `json:"location_id"`
	Status     int `json:"status"`
}

func (ProgramLocationTb) TableName() string { return "program_location_tb" }

func (ProgramLocationTb) ColumnId() string         { return "id" }
func (ProgramLocationTb) ColumnProgramId() string  { return "program_id" }
func (ProgramLocationTb) ColumnLocationId() string { return "location_id" }
func (ProgramLocationTb) ColumnStatus() string     { return "status" }
