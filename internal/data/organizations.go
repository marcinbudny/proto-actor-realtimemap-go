package data

type Organization struct {
	id        string
	name      string
	geofences []*CircularGeofence
}

var AllOrganizations = map[string]*Organization{
	"0006": {
		id:   "006",
		name: "Oy Pohjolan Liikenne Ab",
	},
	"0012": {
		id:        "0012",
		name:      "Helsingin Bussiliikenne Oy",
		geofences: []*CircularGeofence{Airport, KallioDistrict, RailwaySquare},
	},
	"0017": {
		id:        "0017",
		name:      "Tammelundin Liikenne Oy",
		geofences: []*CircularGeofence{LaajasaloIsland},
	},
	"0018": {
		id:        "0018",
		name:      "Pohjolan Kaupunkiliikenne Oy",
		geofences: []*CircularGeofence{KallioDistrict, LauttasaariIsland, RailwaySquare},
	},
	"0020": {
		id:   "0020",
		name: "Bus Travel Åbergin Linja Oy",
	},
	"0021": {
		id:   "0021",
		name: "Bus Travel Oy Reissu Ruoti",
	},
	"0022": {
		id:        "0022",
		name:      "Nobina Finland Oy",
		geofences: []*CircularGeofence{Airport, KallioDistrict, LaajasaloIsland},
	},
	"0030": {
		id:        "0030",
		name:      "Savonlinja Oy",
		geofences: []*CircularGeofence{Airport, Downtown},
	},
	"0036": {
		id:   "0036",
		name: "Nurmijärven Linja Oy",
	},
	"0040": {
		id:   "0040",
		name: "HKL-Raitioliikenne",
	},
	"0045": {
		id:   "0045",
		name: "Transdev Vantaa Oy",
	},
	"0047": {
		id:   "0047",
		name: "Taksikuljetus Oy",
	},
	"0050": {
		id:   "0050",
		name: "HKL-Metroliikenne",
	},
	"0051": {
		id:   "0051",
		name: "Korsisaari Oy",
	},
	"0054": {
		id:   "0054",
		name: "V-S Bussipalvelut Oy",
	},
	"0055": {
		id:   "0055",
		name: "Transdev Helsinki Oy",
	},
	"0058": {
		id:   "0058",
		name: "Koillisen Liikennepalvelut Oy",
	},
	"0060": {
		id:   "0060",
		name: "Suomenlinnan Liikenne Oy",
	},
	"0059": {
		id:   "0059",
		name: "Tilausliikenne Nikkanen Oy",
	},
	"0089": {
		id:   "0089",
		name: "Metropolia",
	},
	"0090": {
		id:   "0090",
		name: "VR Oy",
	},
	"0195": {
		id:   "0195",
		name: "Siuntio1",
	},
}
