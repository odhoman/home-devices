package common

type HomdeDeviceResponse struct {
	ID         string `json:"id"`
	MAC        string `json:"mac"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	HomeID     string `json:"homeId"`
	CreatedAt  int64  `json:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt"`
}
