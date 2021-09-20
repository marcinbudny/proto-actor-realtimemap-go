package contract

type Organization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationDetails struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Geofences []*Geofence `json:"geofences"`
}

type Geofence struct {
	Name           string   `json:"name"`
	Longitude      float64  `json:"longitude"`
	Latitude       float64  `json:"latitude"`
	RadiusInMeters float64  `json:"radiusInMeters"`
	VehiclesInZone []string `json:"vehiclesInZone"`
}
