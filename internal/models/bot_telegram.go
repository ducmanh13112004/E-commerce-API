package models

type ForwardMessage struct {
	Time        string `json:"time" validate:"required"`
	ServiceName string `json:"service_name" validate:"required"`
	Status      int    `json:"status" validate:"required"`
	Message     string `json:"message" validate:"required"`
}
