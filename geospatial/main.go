package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/morikuni/go-geoplot"

	"geospatial/customMap"
	"geospatial/hub"
	"geospatial/node"
	"geospatial/planner"
)

type Config struct {
	Planner PlannerConfig
	Hub     HubConfig
	Node    NodeConfig
}

type PlannerConfig struct {
	Host             string
	Port             int
	AllocationPolicy string
	SimTime          int
}

type HubConfig struct {
	Host           string
	PortRangeStart int
	Hub1Lat        float64
	Hub2Long       float64
	ChargingSpots  int
}

type NodeConfig struct {
	Host           string
	PortRangeStart int
	NumNodes       int
}

func main() {
	// Read configuration file
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	plannerListenPort := strconv.Itoa(config.Planner.Port)

	// Create new Planner
	planner := planner.NewPlanner(config.Planner.AllocationPolicy)

	// create new Hubs
	// Harney Science Center coordinates
	// centerLat := config.Hub.Hub1Lat
	// centerLong := config.Hub.Hub2Long
	centerLat := 37.77695
	centerLong := -122.45117
	hub1 := hub.NewHub(config.Hub.Host, strconv.Itoa(config.Hub.PortRangeStart), config.Planner.AllocationPolicy, centerLat, centerLong, config.Hub.ChargingSpots)

	hub2Lat := centerLat - 0.05
	hub2Long := centerLong + 0.05
	hub2 := hub.NewHub(config.Hub.Host, strconv.Itoa(config.Hub.PortRangeStart+1), config.Planner.AllocationPolicy, hub2Lat, hub2Long, config.Hub.ChargingSpots)

	planner.RegisterHub(hub1)
	planner.RegisterHub(hub2)

	hubCollection := planner.HubMap

	// Add hub to map
	hubGeoplotCoords := &geoplot.LatLng{
		Latitude:  centerLat,
		Longitude: centerLong,
	}
	Area := &geoplot.Area{
		From: hubGeoplotCoords.Offset(-0.1, -0.1),
		To:   hubGeoplotCoords.Offset(0.2, 0.2),
	}

	// Change to general map center
	hubMap := customMap.NewCustomMap(hubGeoplotCoords, 7, Area)
	customMap.InitHub(hubCollection, hubMap)

	// Start Map Server
	go customMap.StartMapServer(hubMap)
	customMap.UpdateMap(hubCollection, hubMap, planner)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				customMap.UpdateMap(hubCollection, hubMap, planner)
			}
		}
	}()

	// Start Hubs
	for _, hub := range planner.HubMap {
		go hub.Listen()
	}

	// Start Planner
	go planner.StartListening(plannerListenPort)

	// Setup Nodes
	plannerHost := "localhost"
	plannerPort := plannerListenPort
	nodeListenPort := 10000
	var nodeArray []*node.NodeStruct

	idCounter := 100
	// Create new Nodes
	for i := 0; i < config.Node.NumNodes; i++ {
		// Node Starting Coordinates
		// Alternate between hub1 and hub2
		nodeLat := centerLat
		nodeLong := centerLong
		if i%2 == 1 {
			nodeLat = hub2Lat
			nodeLong = hub2Long
		}
		n, err := node.NewNode(plannerHost+":"+plannerPort, strconv.Itoa(nodeListenPort), nodeLat, nodeLong)
		if err != nil {
			fmt.Println("Error creating node", i)
			return
		}
		n.NodeId = strconv.Itoa(idCounter)
		fmt.Print("Created node ", i, " with ID: ", n.NodeId, " and listening: ", n.NodeAddr, "\n")
		nodeArray = append(nodeArray, n)
		nodeListenPort++
		idCounter++
	}

	fmt.Println("Nodes created: ", len(nodeArray))

	// Start Nodes

	for i := 0; i < config.Node.NumNodes; i++ {
		fmt.Println("Starting node", i)
		go nodeArray[i].ActivateNode()
	}

	// Signals to Stop / Interrupt main simulation
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	simTime := time.After(time.Duration(config.Planner.SimTime) * time.Second)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		fmt.Println("Shutting down...")

		done <- true
	}()

	go func() {
		<-simTime
		fmt.Println()
		fmt.Println("Simulation Time Reached")

		done <- true
	}()

	// Wait for interrupt signal
	fmt.Println("Simulation is running, press Ctrl+C to stop.")
	<-done

	time.Sleep(1 * time.Second)
	fmt.Println()

	fmt.Printf("--------------------------------------------------\n\n")
	fmt.Printf("------------------DOWNTIME STATS------------------\n")
	fmt.Println("HUB 1")
	fmt.Printf("Total Downtime: %v\n", hub1.GetHubDowntime())
	avgHub1 := hub1.GetAverageDowntimeStat()
	fmt.Printf("Average Downtime per node: %v\n", avgHub1)

	fmt.Println()

	fmt.Println("HUB 2")
	fmt.Printf("Total Downtime: %v\n", hub2.GetHubDowntime())
	avgHub2 := hub1.GetAverageDowntimeStat()
	fmt.Printf("Average Downtime per node: %v\n", avgHub2)

	fmt.Printf("--------------------------------------------------\n\n")

	fmt.Println("Exited")

}
