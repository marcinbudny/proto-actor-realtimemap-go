package contract

type Position struct {
	VehicleId string  `json:"vehicleId"`
	OrgId     string  `json:"orgId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
	Heading   int32   `json:"heading"`
	Speed     float64 `json:"speed"`
	DoorsOpen bool    `json:"doorsOpen"`
}

type PositionBatch struct {
	Positions []*Position `json:"positions"`
}
