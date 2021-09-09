package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/marcinbudny/realtimemap-go/contract"
	"github.com/philippseith/signalr"
)

type AppHub struct {
	signalr.Hub
	initialized bool
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
}

func (h *AppHub) SendPosition(position *contract.Position) {
	serialized, _ := json.Marshal(position)
	if h.initialized {
		h.Clients().All().Send("position", string(serialized))
	}
}

func serveHub(router *http.ServeMux, ctx context.Context) *AppHub {
	hub := &AppHub{}

	singnalrServer, _ := signalr.NewServer(ctx,
		signalr.UseHub(hub),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false))

	singnalrServer.MapHTTP(router, "/events")

	return hub
}
