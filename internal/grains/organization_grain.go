package grains

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/marcinbudny/realtimemap-go/internal/data"
)

type organizationGrain struct {
	id          string
	name        string
	initialized bool
	cluster     *cluster.Cluster
}

func CreateOrganizationFactory(cluster *cluster.Cluster) func() Organization {
	return func() Organization {
		return &organizationGrain{cluster: cluster}
	}
}

func (o *organizationGrain) Init(id string) {
	o.id = id
}

func (o *organizationGrain) initializeOrgIfNeeded(ctx cluster.GrainContext) {
	if o.initialized {
		return
	}

	if organization, ok := data.AllOrganizations[o.id]; ok {
		o.name = organization.Name
		for _, geofence := range organization.Geofences {
			o.createGeofenceActor(geofence, ctx)
		}
	}

	o.initialized = true
}

func (o *organizationGrain) createGeofenceActor(geofence *data.CircularGeofence, ctx cluster.GrainContext) {
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewGeofenceActor(o.name, geofence, o.cluster)
	})
	ctx.Spawn(props)
}

func (o *organizationGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {
	// TODO: normally this would be peformed in the Init func, but the generated code does not pass context to it
	o.initializeOrgIfNeeded(ctx)

	for _, geofenceActor := range ctx.Children() {
		ctx.Send(geofenceActor, position)
	}

	return &Empty{}, nil
}

func (o *organizationGrain) GetGeofences(*GetGeofencesRequest, cluster.GrainContext) (*GetGeofencesResponse, error) {
	return &GetGeofencesResponse{Geofences: []*GeofenceDetails{}}, nil
}

func (o *organizationGrain) Terminate()                       {}
func (o *organizationGrain) ReceiveDefault(ctx actor.Context) {}
