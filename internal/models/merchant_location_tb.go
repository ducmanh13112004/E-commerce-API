package models

import "time"

type MerchantLocationTb struct {
	LocationId            int        `json:"location_id"`
	MerchantId            string     `json:"merchant_id"`
	StoreName             string     `json:"store_name"`
	FullAddress           string     `json:"full_address"`
	Coordinate            string     `json:"coordinate"`
	LocationIcon          string     `json:"location_icon"`
	Address               string     `json:"address"`
	Ward                  string     `json:"ward"`
	District              string     `json:"district"`
	Province              string     `json:"province"`
	Activate              int        `json:"activate"`
	TCreate               *time.Time `json:"t_create"`
	Type                  string     `json:"type"`
	AddressWifiAreaEnable string     `json:"address_wifi_area_enable"`
	Location              string     `json:"location"`
	MacAddress            string     `json:"mac_address"`
}

func (MerchantLocationTb) TableName() string { return "merchant_location_tb" }

func (MerchantLocationTb) ColumnLocationId() string            { return "location_id" }
func (MerchantLocationTb) ColumnMerchantId() string            { return "merchant_id" }
func (MerchantLocationTb) ColumnStoreName() string             { return "store_name" }
func (MerchantLocationTb) ColumnFullAddress() string           { return "full_address" }
func (MerchantLocationTb) ColumnCoordinate() string            { return "coordinate" }
func (MerchantLocationTb) ColumnLocationIcon() string          { return "location_icon" }
func (MerchantLocationTb) ColumnAddress() string               { return "address" }
func (MerchantLocationTb) ColumnWard() string                  { return "ward" }
func (MerchantLocationTb) ColumnDistrict() string              { return "district" }
func (MerchantLocationTb) ColumnProvince() string              { return "province" }
func (MerchantLocationTb) ColumnActivate() string              { return "activate" }
func (MerchantLocationTb) ColumnTCreate() string               { return "t_create" }
func (MerchantLocationTb) ColumnType() string                  { return "type" }
func (MerchantLocationTb) ColumnAddressWifiAreaEnable() string { return "address_wifi_area_enable" }
func (MerchantLocationTb) ColumnLocation() string              { return "location" }
func (MerchantLocationTb) ColumnMacAddress() string            { return "mac_address" }

type LocationProgram struct {
	MerchantId   string `json:"merchant_id"`
	StoreName    string `json:"store_name"`
	FullAddress  string `json:"full_address"`
	Coordinate   string `json:"coordinate"`
	LocationIcon string `json:"location_icon"`
}
