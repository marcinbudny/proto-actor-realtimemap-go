package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	httpapi "github.com/marcinbudny/realtimemap-go/internal/api"
	"github.com/marcinbudny/realtimemap-go/internal/grains"
	"github.com/marcinbudny/realtimemap-go/internal/ingress"
	"github.com/marcinbudny/realtimemap-go/internal/protocluster"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stopOnSignals(cancel)

	cluster := protocluster.StartNode()
	time.Sleep(2 * time.Second)

	api := httpapi.NewHttpApi(cluster.ActorSystem, ctx)
	httpDone := api.ListenAndServe()

	ingressDone := ingress.ConsumeVehicleEvents(func(event *ingress.Event) {
		position := MapToPosition(event)
		if position != nil {
			vehicleGrainClient := grains.GetVehicleGrainClient(cluster, position.VehicleId)
			vehicleGrainClient.OnPosition(position)
		}
	}, ctx)

	<-ingressDone
	<-httpDone
	// cluster.ActorSystem.EventStream.Unsubscribe(subscription)
	// close(batchingChan)
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
		OrgId:     e.OperatorId,
		Latitude:  *payload.Latitude,
		Longitude: *payload.Longitude,
		Heading:   *payload.Heading,
		Timestamp: (*payload.Timestamp).UnixMilli(),
		Speed:     *payload.Speed,
	}
}

// func MapToHubPosition(p *grains.Position) *contract.Position {
// 	return &contract.Position{
// 		VehicleId: p.VehicleId,
// 		OrgId:     p.OrgId,
// 		Latitude:  p.Latitude,
// 		Longitude: p.Longitude,
// 		Heading:   p.Heading,
// 		Timestamp: p.Timestamp,
// 		Speed:     p.Speed,
// 	}
// }
