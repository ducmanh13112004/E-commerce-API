package models

import "time"

type StoreCallApiTb struct {
	Id        int       `json:"id"`
	Url       string    `json:"url"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	TCreate   time.Time `json:"t_create"`
	TResponse time.Time `json:"t_response"`
	Dt        float64   `json:"dt"`
}

func (StoreCallApiTb) TableName() string     { return "store_call_api_tb" }
func (StoreCallApiTb) ColumnId() string      { return "id" }
func (StoreCallApiTb) ColumnUrl() string     { return "url" }
func (StoreCallApiTb) ColumnInput() string   { return "input" }
func (StoreCallApiTb) ColumnOutput() string  { return "output" }
func (StoreCallApiTb) ColumnTCreate() string { return "t_create" }
