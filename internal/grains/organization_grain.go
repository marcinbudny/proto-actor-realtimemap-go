package grains

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

type organizationGrain struct {
	id      string
	cluster *cluster.Cluster
}

func CreateOrganizationFactory(cluster *cluster.Cluster) func() Organization {
	return func() Organization {
		return &organizationGrain{cluster: cluster}
	}
}

func (o *organizationGrain) Init(id string) {
	o.id = id
}

func (o *organizationGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Empty, error) {
	return &Empty{}, nil
}

func (o *organizationGrain) Terminate()                       {}
func (o *organizationGrain) ReceiveDefault(ctx actor.Context) {}
