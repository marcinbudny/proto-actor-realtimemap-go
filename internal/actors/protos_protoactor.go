
package actors

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var rootContext = actor.EmptyRootContext
	
var xVehicleFactory func() Vehicle

// VehicleFactory produces a Vehicle
func VehicleFactory(factory func() Vehicle) {
	xVehicleFactory = factory
}

// GetVehicleGrain instantiates a new VehicleGrain with given ID
func GetVehicleGrain(id string) *VehicleGrain {
	return &VehicleGrain{ID: id}
}

// Vehicle interfaces the services available to the Vehicle
type Vehicle interface {
	Init(id string)
	Terminate()
		
	OnPosition(*Position, cluster.GrainContext) (*Noop, error)
		
}

// VehicleGrain holds the base data for the VehicleGrain
type VehicleGrain struct {
	ID string
}
	
// OnPosition requests the execution on to the cluster using default options
func (g *VehicleGrain) OnPosition(r *Position) (*Noop, error) {
	return g.OnPositionWithOpts(r, cluster.DefaultGrainCallOptions())
}

// OnPositionWithOpts requests the execution on to the cluster
func (g *VehicleGrain) OnPositionWithOpts(r *Position, opts *cluster.GrainCallOptions) (*Noop, error) {
	fun := func() (*Noop, error) {
			pid, statusCode := cluster.Get(g.ID, "Vehicle")
			if statusCode != remote.ResponseStatusCodeOK && statusCode != remote.ResponseStatusCodePROCESSNAMEALREADYEXIST {
				return nil, fmt.Errorf("get PID failed with StatusCode: %v", statusCode)
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
			response, err := rootContext.RequestFuture(pid, request, opts.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &Noop{}
				err = proto.Unmarshal(msg.MessageData, result)
				if err != nil {
					return nil, err
				}
				return result, nil
			case *cluster.GrainErrorResponse:
				return nil, errors.New(msg.Err)
			default:
				return nil, errors.New("unknown response")
			}
		}
	
	var res *Noop
	var err error
	for i := 0; i < opts.RetryCount; i++ {
		res, err = fun()
		if err == nil || err.Error() != "future: timeout" {
			return res, err
		} else if opts.RetryAction != nil {
				opts.RetryAction(i)
		}
	}
	return nil, err
}

// OnPositionChan allows to use a channel to execute the method using default options
func (g *VehicleGrain) OnPositionChan(r *Position) (<-chan *Noop, <-chan error) {
	return g.OnPositionChanWithOpts(r, cluster.DefaultGrainCallOptions())
}

// OnPositionChanWithOpts allows to use a channel to execute the method
func (g *VehicleGrain) OnPositionChanWithOpts(r *Position, opts *cluster.GrainCallOptions) (<-chan *Noop, <-chan error) {
	c := make(chan *Noop)
	e := make(chan error)
	go func() {
		res, err := g.OnPositionWithOpts(r, opts)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}
	

// VehicleActor represents the actor structure
type VehicleActor struct {
	inner Vehicle
	Timeout *time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *VehicleActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xVehicleFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:]) // skip "remote$"
		if a.Timeout != nil {
			ctx.SetReceiveTimeout(*a.Timeout)
		}
	case *actor.ReceiveTimeout:
		a.inner.Terminate()
		ctx.Self().Poison()

	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
			
		case 0:
			req := &Position{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.OnPosition(req, ctx)
			if err == nil {
				bytes, errMarshal := proto.Marshal(r0)
				if errMarshal != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", errMarshal)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}
		
		}
	default:
		log.Printf("Unknown message %v", msg)
	}
}

	



