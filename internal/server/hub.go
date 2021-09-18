package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/AsynkronIT/protoactor-go/actor"
	kitlog "github.com/go-kit/log"
	"github.com/marcinbudny/realtimemap-go/internal/grains"
	"github.com/philippseith/signalr"
)

type AppHub struct {
	signalr.Hub
	initialized bool
	actorSystem *actor.ActorSystem
}

func (h *AppHub) Initialize(ctx signalr.HubContext) {
	h.Hub.Initialize(ctx)
	// Initialize will be called on first connection to the hub
	// if position is sent before that, the HubContext is nil and application crashes
	// TODO: possible bug in singlar implementation?
	h.initialized = true
}

func (h *AppHub) OnConnected(connectionID string) {
	fmt.Printf("%s connected to hub\n", connectionID)
	h.Clients().Caller().Send("event", "{\"hello\": \"world\"}")
}

func (h *AppHub) OnDisconnected(connectionID string) {
	fmt.Printf("%s disconnected from hub\n", connectionID)
	if item, loaded := h.Items().LoadAndDelete("viewport"); loaded {
		viewportPID, _ := item.(*actor.PID)
		h.actorSystem.Root.Stop(viewportPID)
	}
}

func (h *AppHub) SendPositionBatch(connectionID string, batch *grains.PositionBatch) {
	serialized, _ := json.Marshal(batch)
	if h.initialized {
		h.Clients().Client(connectionID).Send("positions", string(serialized))
	}
}

func (h *AppHub) SendNotification(connectionID string, message string) {
	if h.initialized {
		h.Clients().Client(connectionID).Send("notification", message)
	}
}

func (h *AppHub) SetViewport(swLng float64, swLat float64, neLng float64, neLat float64) {
	var viewportPID *actor.PID

	if item, loaded := h.Items().Load(h.ConnectionID()); loaded {
		viewportPID, _ = item.(*actor.PID)
	} else {
		props := actor.PropsFromProducer(func() actor.Actor {
			return grains.NewViewportActor(h.SendPositionBatch, h.SendNotification)
		})
		// spawn named so that we don't get multiple viewports for same connection id in the case of concurrency issues
		viewportPID, _ = h.actorSystem.Root.SpawnNamed(props, h.ConnectionID())
		h.Items().Store("viewport", viewportPID)
	}

	h.actorSystem.Root.Send(viewportPID, &grains.UpdateViewport{
		Viewport: &grains.Viewport{
			SouthWest: &grains.GeoPoint{Longitude: swLng, Latitude: swLat},
			NorthEast: &grains.GeoPoint{Longitude: neLng, Latitude: neLat},
		},
	})
}

func serveHub(router *http.ServeMux, actorSystem *actor.ActorSystem, ctx context.Context) *AppHub {
	hub := &AppHub{actorSystem: actorSystem}

	singnalrServer, _ := signalr.NewServer(ctx,
		signalr.UseHub(hub),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false))

	singnalrServer.MapHTTP(router, "/events")

	return hub
}
