package grains

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/eventstream"
)

const BatchSize = 10

type SendPositions func(*PositionBatch)

type viewportActor struct {
	connectionId  string
	viewport      Viewport
	batch         []*Position
	sendPositions SendPositions
	subscription  *eventstream.Subscription
}

func NewViewportActor(sendPositions SendPositions) *viewportActor {
	return &viewportActor{sendPositions: sendPositions}
}

func (v *viewportActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		v.batch = make([]*Position, 0, BatchSize)

		v.subscription = ctx.ActorSystem().EventStream.Subscribe(func(event interface{}) {
			if pos, ok := event.(*Position); ok {
				ctx.Send(ctx.Self(), pos)
			}
		})

	case *Position:
		if msg.Latitude < v.viewport.SouthWest.Latitude ||
			msg.Latitude > v.viewport.NorthEast.Latitude ||
			msg.Longitude < v.viewport.SouthWest.Longitude ||
			msg.Longitude > v.viewport.NorthEast.Longitude {
			return
		}

		v.batch = append(v.batch, msg)
		if len(v.batch) >= BatchSize {
			v.sendPositions(&PositionBatch{Positions: v.batch})
			v.batch = v.batch[:0]
		}

	case *UpdateViewport:
		v.viewport = *msg.Viewport
		fmt.Printf("Viewport for connection %s is now %+v\n", v.connectionId, v.viewport)

	case *actor.Stopping:
		fmt.Printf("Stopping viewport for connection %s\n", v.connectionId)
		ctx.ActorSystem().EventStream.Unsubscribe(v.subscription)
		v.subscription = nil
	}
}
