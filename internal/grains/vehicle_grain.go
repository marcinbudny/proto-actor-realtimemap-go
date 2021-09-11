package grains

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

type VehicleGrain struct {
	id string
}

func (v *VehicleGrain) Init(id string) {
	v.id = id
}

func (v *VehicleGrain) OnPosition(position *Position, ctx cluster.GrainContext) (*Noop, error) {

	fmt.Printf("Vehicle %s received %+v\n", v.id, position)

	return &Noop{}, nil
}

func (v *VehicleGrain) Terminate()                       {}
func (v *VehicleGrain) ReceiveDefault(ctx actor.Context) {}
