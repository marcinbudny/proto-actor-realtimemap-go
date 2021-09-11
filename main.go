package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/marcinbudny/realtimemap-go/contract"
	httpapi "github.com/marcinbudny/realtimemap-go/internal/api"
	"github.com/marcinbudny/realtimemap-go/internal/grains"
	"github.com/marcinbudny/realtimemap-go/internal/ingress"
	"github.com/marcinbudny/realtimemap-go/internal/protocluster"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stopOnSignals(cancel)

	cluster, _ := protocluster.StartNode()
	api := httpapi.NewHttpApi(ctx)
	httpDone := api.ListenAndServe()

	ingressDone := ingress.ConsumeVehicleEvents(func(event *ingress.Event) {
		hubPosition := MapToHubPosition(event)
		if hubPosition != nil {
			api.Hub.SendPosition(hubPosition)
		}
		position := MapToPosition(event)
		if position != nil {
			vehicleGrainClient := grains.GetVehicleGrainClient(cluster, position.VehicleId)
			vehicleGrainClient.OnPosition(position)
		}
	}, ctx)

	<-ingressDone
	<-httpDone
	cluster.Shutdown(true)
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

func MapToPosition(e *ingress.Event) *grains.Position {
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

	return &grains.Position{
		VehicleId: e.VehicleId,
		Latitude:  *payload.Latitude,
		Longitude: *payload.Longitude,
		Heading:   *payload.Heading,
		Timestamp: (*payload.Timestamp).UnixMilli(),
		Speed:     *payload.Speed,
	}
}

func MapToHubPosition(e *ingress.Event) *contract.Position {
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
