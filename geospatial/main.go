package main

import (
	"os"
	"fmt"
	"net"
	"time"
	
	"geospatial/messages"
	"geospatial/hub"
	"geospatial/customMap"
	"geospatial/planner"

	"github.com/morikuni/go-geoplot"
)



func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./hub <listen-port>")
	}

	listenPort := os.Args[1]

	// Listen for incoming connections
	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		fmt.Println("Error starting hub server:", err)
		return
	}

	defer listener.Close()

	// Create new Planner
	planner := planner.NewPlanner()

	// create new Hubs
	// Harney Science Center coordinates
	centerLat := 37.77695
	centerLong := -122.45117
	hubStation := hub.NewHub("orion02","9001", centerLat, centerLong)

	hub2 := hub.NewHub("orion02", "9002", centerLat - 0.05, centerLong + 0.05)

	planner.RegisterHub(hubStation)
	planner.RegisterHub(hub2)

	hubCollection := planner.HubMap

	// Add hub to map
	hubGeoplotCoords := &geoplot.LatLng {
		Latitude: centerLat,
		Longitude: centerLong,
	}
	Area:= &geoplot.Area{
		From: hubGeoplotCoords.Offset(-0.1, -0.1),
		To:   hubGeoplotCoords.Offset(0.2, 0.2),
	}

	// Change to general map center
	hubMap := customMap.NewCustomMap(hubGeoplotCoords, 7, Area)
	customMap.InitHub(hubCollection, hubMap)

	go customMap.StartMapServer(hubMap)
	customMap.UpdateMap(hubCollection, hubMap)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				customMap.UpdateMap(hubCollection, hubMap)
			}
		}
	}()

	for _, hub := range planner.HubMap {
		go hub.Listen()
	} 

	for {
		// Accept incoming connections to Planner
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Message Handler
		msgHandler := messages.NewMessageHandler(conn)
		//go hub.HandleMessage(conn, msgHandler, hubStation)
		go planner.HandleMessage(conn, msgHandler)
	}

}