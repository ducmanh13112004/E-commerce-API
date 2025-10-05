package models

type UserInfo struct {
	CustomerId  string
	PhoneNb     string
	AccessToken string
	AppVersion  string
	TokenWebkit string
	CustomerIp  string
}
type ResultHealthCheck struct {
	RequestID string            `json:"request_id"`
	Status    string            `json:"status"`
	Detail    DetailHealthCheck `json:"detail"`
}
type DetailHealthCheck struct {
	RAMIndex        string `json:"ram_index"`
	CPUIndex        string `json:"cpu_index"`
	DatabaseIndex   string `json:"database_index"`
	RedisIndex      string `json:"redis_index"`
	KafkaIndex      string `json:"kafka_index"`
	ThirdPartyIndex string `json:"3rd_index"`
}
type APICustomerInfo struct {
	CustomerId            int    `json:"customerId"`
	Phone                 string `json:"phone"`
	FullName              string `json:"fullName"`
	Address               string `json:"address"`
	Language              string `json:"language"`
	Avatar                string `json:"avatar"`
	Email                 string `json:"email"`
	Birthday              string `json:"birthday"`
	Gender                string `json:"gender"`
	FirstLogin            string `json:"firstLogin"`
	LastLogin             string `json:"lastLogin"`
	Version               string `json:"version"`
	SessionDevicePlatform string `json:"sessionDevicePlatform"`
	SessionDeviceId       string `json:"sessionDeviceId"`
	SessionDeviceToken    string `json:"sessionDeviceToken"`
	SessionDeviceName     string `json:"sessionDeviceName"`
	// LoyaltyStatus         *int  `json:"loyaltyStatus"`
}
