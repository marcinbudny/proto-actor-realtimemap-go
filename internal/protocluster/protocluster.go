package protocluster

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/automanaged"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/marcinbudny/realtimemap-go/internal/grains"
)

func StartNode() (*cluster.Cluster, *actor.ActorSystem) {
	system := actor.NewActorSystem()

	vehicleKind := cluster.NewKind("Vehicle", actor.PropsFromProducer((func() actor.Actor {
		return &grains.VehicleActor{}
	})))

	provider := automanaged.New()
	config := remote.Configure("localhost", 0)

	clusterConfig := cluster.Configure("my-cluster", provider, config, vehicleKind)
	cluster := cluster.New(system, clusterConfig)

	grains.VehicleFactory(func() grains.Vehicle { return &grains.VehicleGrain{} })

	cluster.Start()

	return cluster, system
}
