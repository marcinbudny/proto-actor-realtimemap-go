package contract

type Position struct {
	VehicleId string  `json:"vehicleId"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
	Heading   int32   `json:"heading"`
	Speed     float32 `json:"speed"`
	DoorsOpen bool    `json:"doorsOpen"`
}
