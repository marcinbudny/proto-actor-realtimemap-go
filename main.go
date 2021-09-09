package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/marcinbudny/realtimemap-go/contract"
	httpapi "github.com/marcinbudny/realtimemap-go/internal/api"
	"github.com/marcinbudny/realtimemap-go/internal/ingress"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stopOnSignals(cancel)

	api := httpapi.NewHttpApi(ctx)
	httpDone := api.ListenAndServe()

	ingressDone := ingress.ConsumeVehicleEvents(func(event *ingress.Event) {
		position := MapToPosition(event)
		if position != nil {
			api.Hub.SendPosition(position)
		}
	}, ctx)

	<-ingressDone
	<-httpDone
}

func stopOnSignals(cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println("*** STOPPING ***")
		cancel()
	}()
}

func MapToPosition(e *ingress.Event) *contract.Position {
	var payload *ingress.Payload

	if e.VehiclePosition != nil {
		payload = e.VehiclePosition
	} else if e.DoorOpen != nil {
		payload = e.DoorOpen
	} else if e.DoorClosed != nil {
		payload = e.DoorClosed
	} else {
		return nil
	}

	if !payload.HasValidPosition() {
		return nil
	}

	return &contract.Position{
		VehicleId: e.VehicleId,
		Latitude:  *payload.Latitude,
		Longitude: *payload.Longitude,
		Heading:   *payload.Heading,
		Timestamp: (*payload.Timestamp).UnixMilli(),
		Speed:     *payload.Speed,
	}
}
