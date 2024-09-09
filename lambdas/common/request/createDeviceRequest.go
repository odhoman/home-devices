package common

type CreateDeviceRequest struct {
	MAC    string `json:"mac" validate:"required,min=12,max=17"`
	Name   string `json:"name" validate:"required,min=3,max=50"`
	Type   string `json:"type" validate:"required,min=3,max=20"`
	HomeID string `json:"homeId" validate:"required,min=5,max=30"`
}
