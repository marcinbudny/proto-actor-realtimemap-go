// Package grains is generated by protoactor-go/protoc-gen-gograin@0.1.0
package grains

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/remote"
	logmod "github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/proto"
)

var (
	plog = logmod.New(logmod.InfoLevel, "[GRAIN]")
	_    = proto.Marshal
	_    = fmt.Errorf
	_    = math.Inf
)

// SetLogLevel sets the log level.
func SetLogLevel(level logmod.Level) {
	plog.SetLevel(level)
}

var xVehicleFactory func() Vehicle

// VehicleFactory produces a Vehicle
func VehicleFactory(factory func() Vehicle) {
	xVehicleFactory = factory
}

// GetVehicleGrainClient instantiates a new VehicleGrainClient with given ID
func GetVehicleGrainClient(c *cluster.Cluster, id string) *VehicleGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &VehicleGrainClient{ID: id, cluster: c}
}

// Vehicle interfaces the services available to the Vehicle
type Vehicle interface {
	Init(id string)
	Terminate()
	ReceiveDefault(ctx actor.Context)
	OnPosition(*Position, cluster.GrainContext) (*Noop, error)
	
}

// VehicleGrainClient holds the base data for the VehicleGrain
type VehicleGrainClient struct {
	ID      string
	cluster *cluster.Cluster
}

// OnPosition requests the execution on to the cluster with CallOptions
func (g *VehicleGrainClient) OnPosition(r *Position, opts ...*cluster.GrainCallOptions) (*Noop, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Call(g.ID, "Vehicle", reqMsg, opts...)
	if err != nil {
		return nil, err
	}
	switch msg := resp.(type) {
	case *cluster.GrainResponse:
		result := &Noop{}
		err = proto.Unmarshal(msg.MessageData, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case *cluster.GrainErrorResponse:
		if msg.Code == remote.ResponseStatusCodeDeadLetter.ToInt32() {
			return nil, remote.ErrDeadLetter
		}
		return nil, errors.New(msg.Err)
	default:
		return nil, errors.New("unknown response")
	}
}


// VehicleActor represents the actor structure
type VehicleActor struct {
	inner   Vehicle
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *VehicleActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *cluster.ClusterInit:
		a.inner = xVehicleFactory()
		a.inner.Init(msg.ID)
		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}

	case *actor.ReceiveTimeout:
		a.inner.Terminate()
		ctx.Poison(ctx.Self())

	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
		case 0:
			req := &Position{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				plog.Error("OnPosition(Position) proto.Unmarshal failed.", logmod.Error(err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			r0, err := a.inner.OnPosition(req, ctx)
			if err != nil {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			bytes, err := proto.Marshal(r0)
			if err != nil {
				plog.Error("OnPosition(Position) proto.Marshal failed", logmod.Error(err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			resp := &cluster.GrainResponse{MessageData: bytes}
			ctx.Respond(resp)
		
		}
	default:
		a.inner.ReceiveDefault(ctx)
	}
}