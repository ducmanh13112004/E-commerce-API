package models

type ProgramCustomerVipTb struct {
	CustomerId        string `json:"customer_id"`
	ProgramId         string `json:"program_id"`
	TimesEachCustomer int    `json:"times_each_customer"`
}

func (ProgramCustomerVipTb) TableName() string { return "program_customer_vip_tb" }

func (ProgramCustomerVipTb) ColumnCustomerId() string        { return "customer_id" }
func (ProgramCustomerVipTb) ColumnProgramId() string         { return "program_id" }
func (ProgramCustomerVipTb) ColumnTimesEachCustomer() string { return "times_each_customer" }
