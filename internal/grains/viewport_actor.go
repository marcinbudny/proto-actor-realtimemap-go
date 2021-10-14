package grains

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/eventstream"
)

const BatchSize = 10

type SendPositions func(connectionID string, positions *PositionBatch)
type SendNotification func(connectionID string, message string)

type viewportActor struct {
	connectionID     string
	viewport         Viewport
	batch            []*Position
	sendPositions    SendPositions
	sendNotification SendNotification
	subscription     *eventstream.Subscription
}

func NewViewportActor(connectionID string, sendPositions SendPositions, sendNotification SendNotification) *viewportActor {
	return &viewportActor{
		connectionID:     connectionID,
		sendPositions:    sendPositions,
		sendNotification: sendNotification,
	}
}

func (v *viewportActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		v.batch = make([]*Position, 0, BatchSize)

		v.subscription = ctx.ActorSystem().EventStream.Subscribe(func(event interface{}) {
			// do not modify state in the callback to avoid concurrency issues, let the message pass through mailbox
			switch event.(type) {
			case *Position:
				ctx.Send(ctx.Self(), event)
			case *Notification:
				ctx.Send(ctx.Self(), event)
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
			v.sendPositions(v.connectionID, &PositionBatch{Positions: v.batch})
			v.batch = v.batch[:0]
		}

	case *Notification:
		v.sendNotification(v.connectionID, msg.Message)

	case *UpdateViewport:
		v.viewport = *msg.Viewport
		fmt.Printf("Viewport for connection %s is now %+v\n", v.connectionID, v.viewport)

	case *actor.Stopping:
		fmt.Printf("Stopping viewport for connection %s\n", v.connectionID)
		ctx.ActorSystem().EventStream.Unsubscribe(v.subscription)
		v.subscription = nil
	}
}
