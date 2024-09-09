package common

type UpdateDeviceRequest struct {
	MAC    string `json:"mac" validate:"omitempty,min=12,max=17"`
	Name   string `json:"name" validate:"omitempty,min=3,max=50"`
	Type   string `json:"type" validate:"omitempty,min=3,max=20"`
	HomeID string `json:"homeId" validate:"omitempty,min=5,max=30"`
}
