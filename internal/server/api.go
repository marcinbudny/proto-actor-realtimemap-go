package server

import (
	"net/http"
	"sort"

	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/labstack/echo"
	"github.com/marcinbudny/realtimemap-go/contract"
	"github.com/marcinbudny/realtimemap-go/internal/data"
)

func serveApi(e *echo.Echo, cluster *cluster.Cluster) {

	e.GET("/api/organization", func(c echo.Context) error {
		result := make([]*contract.Organization, 0)

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
}
