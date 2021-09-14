package grains

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

type vehicleGrain struct {
	id      string
	cluster *cluster.Cluster
}

func CreateVehicleFactory(cluster *cluster.Cluster) func() Vehicle {
	return func() Vehicle {
		return &vehicleGrain{cluster: cluster}
	}
}

func (v *vehicleGrain) Init(id string) {
	v.id = id
}

func (v *vehicleGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {

	orgClient := GetOrganizationGrainClient(v.cluster, position.OrgId)
	orgClient.OnPosition(position)

	// the Go version of proto.actor does not support broadcasting among cluster members yet
	// TODO: fix this code once it does
	v.cluster.ActorSystem.EventStream.Publish(position)

	return &Empty{}, nil
}

func (v *vehicleGrain) GetPositionsHistory(*GetPositionsHistoryRequest, cluster.GrainContext) (*PositionBatch, error) {
	return &PositionBatch{}, nil
}

func (v *vehicleGrain) Terminate()                       {}
func (v *vehicleGrain) ReceiveDefault(ctx actor.Context) {}
