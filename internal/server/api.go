package server

import (
	"fmt"
	"net/http"
	"sort"

	echo "github.com/labstack/echo/v4"

	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/marcinbudny/realtimemap-go/contract"
	"github.com/marcinbudny/realtimemap-go/internal/data"
	"github.com/marcinbudny/realtimemap-go/internal/grains"
)

func serveApi(e *echo.Echo, cluster *cluster.Cluster) {

	e.GET("/api/organization", func(c echo.Context) error {
		result := make([]*contract.Organization, 0, len(data.AllOrganizations))

		for _, org := range data.AllOrganizations {
			if len(org.Geofences) > 0 {
				result = append(result, &contract.Organization{
					Id:   org.Id,
					Name: org.Name,
				})
			}
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})

		return c.JSON(http.StatusOK, result)
	})

	e.GET("/api/organization/:id", func(c echo.Context) error {
		var id string
		if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
			c.String(http.StatusBadRequest, "Unable to get id from the request")
		}

		if org, ok := data.AllOrganizations[id]; ok {

			orgClient := grains.GetOrganizationGrainClient(cluster, id)
			if grainResponse, err := orgClient.GetGeofences(&grains.GetGeofencesRequest{}); err == nil {

				return c.JSON(http.StatusOK, mapOrganization(org, grainResponse))

			} else {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to call grain for organization %s: %v", id, err))
			}
		} else {
			return c.String(http.StatusNotFound, fmt.Sprintf("Organization %s not found", id))
		}
	})
}

func mapOrganization(org *data.Organization, grainResponse *grains.GetGeofencesResponse) *contract.OrganizationDetails {
	geofences := make([]*contract.Geofence, 0, len(grainResponse.Geofences))

	for _, grainGeofence := range grainResponse.Geofences {
		geofences = append(geofences, mapGeofence(grainGeofence))
	}

	sort.Slice(geofences, func(i, j int) bool {
		return geofences[i].Name < geofences[j].Name
	})

	return &contract.OrganizationDetails{
		Id:        org.Id,
		Name:      org.Name,
		Geofences: geofences,
	}
}

func mapGeofence(grainGeofence *grains.GeofenceDetails) *contract.Geofence {
	vehicles := make([]string, len(grainGeofence.VehiclesInZone))
	copy(vehicles, grainGeofence.VehiclesInZone)
	sort.Strings(vehicles)

	return &contract.Geofence{
		Name:           grainGeofence.Name,
		Longitude:      grainGeofence.Longitude,
		Latitude:       grainGeofence.Latitude,
		RadiusInMeters: grainGeofence.RadiusInMeters,
		VehiclesInZone: vehicles,
	}
}
