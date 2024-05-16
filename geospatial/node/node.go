package main

import (
	"fmt"
	"os"
	"geospatial/messages"
	"net"
    "time"
	"github.com/golang/protobuf/ptypes"
	"strconv"
	"sync"
	"log"
	"errors"
	"math/rand"
	"math"
)

type nodeStruct struct {
	Mutex sync.Mutex
	plannerAddr string
	hubAddr string
	nodeAddr string
	nodeId	string
	
	lat float64
	long float64
	
	battery_level int32
	seeking_charge bool
	charging bool
	
	hub_set bool
	hub_lat float64
	hub_long float64
	at_hub bool
}

func sendRecon(node *nodeStruct) (bool, string) {
	// connect to hub
	conn, err := net.Dial("tcp", node.plannerAddr)
	if err != nil {
		log.Printf("Could not connect to planner:", node.plannerAddr)
		return false, "false"
	}
	defer conn.Close()

	msgHandler := messages.NewMessageHandler(conn)
	defer msgHandler.Close()

	// Send Recon Request
	coordMsg := messages.Geocoordinate{Latitude: node.lat, Longitude: node.long}
	batteryMsg := messages.Battery{BatteryLevel: node.battery_level, SeekingCharge: node.seeking_charge}

	reconMsg := messages.ReconRequest{NodeAddr: node.nodeAddr, NodeId: node.nodeId, Coordinates: &coordMsg, BatteryState: &batteryMsg}
	wrapper := &messages.Wrapper {
		Msg: &messages.Wrapper_ReconReq{ReconReq: &reconMsg},
	}
	msgHandler.Send(wrapper)

	// Wait for Recon Response
	wrapper, _ = msgHandler.Receive()

	msg := wrapper.GetReconRes()
	//fmt.Println(msg)
	if msg.Accept {
		hubHost := msg.HubHost
		hubPort := msg.HubPort
		node.hubAddr = hubHost + ":" + hubPort
		
		return true, node.hubAddr
	} else {
		log.Println("Node could not recon planner")
		return false, "false"
	}

	return false, "false"	
}

// Sends join message to hub
func sendJoin(node *nodeStruct) (bool, string){
	// connect to hub
	conn, err := net.Dial("tcp", node.hubAddr)
	if err != nil {
		log.Printf("Could not connect to hub:", node.hubAddr)
		return false, "false"
	}
	defer conn.Close()

	msgHandler := messages.NewMessageHandler(conn)
	defer msgHandler.Close()

	// Send Join Request
	coordMsg := messages.Geocoordinate{Latitude: node.lat, Longitude: node.long}
	batteryMsg := messages.Battery{BatteryLevel: node.battery_level, SeekingCharge: node.seeking_charge}

	joinMsg := messages.JoinRequest{NodeAddr: node.nodeAddr, NodeId: node.nodeId, Coordinates: &coordMsg, BatteryState: &batteryMsg}
	wrapper := &messages.Wrapper {
		Msg: &messages.Wrapper_JoinReq{JoinReq: &joinMsg},
	}
	msgHandler.Send(wrapper)

	// Wait for Join Response
	wrapper, _ = msgHandler.Receive()
	msg := wrapper.GetJoinRes()
	if msg.Accept {
		log.Println("Node successfuly joined:", msg.NodeId)
		hubCoordinates := msg.Coordinates
		node.hub_set = true
		node.hub_lat = hubCoordinates.Latitude
		node.hub_long = hubCoordinates.Longitude
		
		return true, msg.NodeId
	} else {
		log.Println("Node could not join hub")
		return false, "false"
	}

	return false, "false"
}

func listen(node *nodeStruct, hbChan chan struct{}) {
	// listen for messages
	listener, err := net.Listen("tcp", node.nodeAddr)
	if err != nil {
		log.Printf("Could not listen on %s\n", node.nodeAddr)
		log.Println(err)
		os.Exit(1)
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			msgHandler := messages.NewMessageHandler(conn)
			go handleMessages(conn, msgHandler, hbChan, node)
		}
	}
}

func handleMessages(conn net.Conn, msgHandler *messages.MessageHandler, hbChan chan struct{}, node *nodeStruct){
	defer msgHandler.Close()
	defer conn.Close()

	for {
		wrapper, _ := msgHandler.Receive()
		
		switch msg := wrapper.Msg.(type) {
		case *messages.Wrapper_KillNode:
			log.Println("Node killed")
			close(hbChan)
			joinStatus, _ := sendJoin(node)
			if joinStatus {
				log.Println("node rejoined, restarting heartbeats")
				sendHearbeats(node, hbChan)
			}

		default:
			log.Printf("Unexpected message type: %T", msg)
			return
		}
	}

}

func sendHearbeats(node *nodeStruct, hbChan <-chan struct{}) error {
	for {
		select {
		case <-hbChan:
			fmt.Println("Node killed, stopping heartbeats")
			return errors.New("node killed")
		default:
			conn, err := net.Dial("tcp", node.hubAddr)
			if err != nil {
				log.Printf("Could not connect to hub: %s\n", node.hubAddr)
			}

			msgHandler := messages.NewMessageHandler(conn)
			defer msgHandler.Close()
			defer conn.Close()

			currentTime := time.Now()
			// convert time to protobuf timestamp
			ts, err := ptypes.TimestampProto(currentTime)
			if err != nil {
				fmt.Println("Error converting time to timestamp: %v\n", err)
				return err
			}

			// UPDATE NODE STATE (Coords & Battery)
			node.Mutex.Lock()
			
			// update battery
			node.battery_level = node.battery_level - 1
			if node.battery_level <= 30 {
				node.seeking_charge = true
			}

			// move node
			movePoint(node)

			node.Mutex.Unlock()

			// Create heartbeat message
			coordMsg := messages.Geocoordinate{Latitude: node.lat, Longitude: node.long}
			batteryMsg := messages.Battery{BatteryLevel: node.battery_level, SeekingCharge: node.seeking_charge}
			hrbtMsg := messages.Heartbeat{NodeId: node.nodeId, State: "Alive", Timestamp: ts, Coordinates: &coordMsg, BatteryState: &batteryMsg}
			wrapper := &messages.Wrapper {
				Msg: &messages.Wrapper_HeartbeatMessage{HeartbeatMessage: &hrbtMsg},
			}
			if err := msgHandler.Send(wrapper); err != nil {
				fmt.Println("Error sending heartbeat message:", err)
				return err
			}
			return nil
		}
	}
}

func generateRandomCoordinate() (float64, float64) {
	// Harney Science Center coordinates
	centerLat := 37.77695
	centerLong := -122.45117

	rand.Seed(time.Now().UnixNano())

	// Generate random adjustments
	deltaLat := (rand.Float64() - 0.5) * (0.05 * 2)
	deltaLong := (rand.Float64() - 0.5) * (0.05 * 2)

	// Apply randomness
	lat := centerLat + deltaLat
	long := centerLong + deltaLong

	return lat, long
}

// Returns a random number between 0 - 100 for battery level
func generateBatteryLevel() (int32) {
	rand.Seed(time.Now().UnixNano())

	// Generate random number between 0 and 100
	randomNumber := rand.Intn(101)

	return int32(randomNumber)
}

// haversine function calculates the distance between two points on the Earth, return kilometers
// given their latitude and longitude in degrees.
func haversine(nodeLat, nodeLong, hubLat, hubLong float64) float64 {
	// convert degrees to radians
	var dlat, dlon, a, c float64
	rad := math.Pi / 180

	hubLat *= rad
	hubLong *= rad
	
	nodeLat *= rad
	nodeLong *= rad

	// Haversine formula
	dlat = nodeLat - hubLat
	dlon = nodeLong - hubLong
	a = math.Pow(math.Sin(dlat/2), 2) + math.Cos(hubLat)*math.Cos(nodeLat)*math.Pow(math.Sin(dlon/2), 2)
	c = 2 * math.Asin(math.Sqrt(a))

	// Radius of Earth in kilometers. Use 3956 for miles
	r := 6371.0

	// calculate and return the distance
	return c * r
}

// movePoint randomly moves the point (lat, lon) a small fixed distance
// up (North), down (South), left (West), or right (East) and returns the new coordinates.
func movePoint(node *nodeStruct) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Define the small distance to move in terms of degrees
	// This is a very rough approximation for demonstration purposes
	const distance = 0.001 // Roughly 111 meters for latitude, varies for longitude

	
	if node.seeking_charge {
		nodeDistance := haversine(node.lat, node.long, node.hub_lat, node.hub_long)
		nodeDistanceMeters := nodeDistance * 1000
		//log.Println("Node Distance to Hub:", nodeDistanceMeters)
		if nodeDistanceMeters <= 30 {
			//log.Println("Node at hub, not moving")
			node.at_hub = true
			node.charging = true

		} else {
			//log.Println("Moving towards hub")
			// Calculate the direction to move in
			latDirection := node.hub_lat - node.lat
			longDirection := node.hub_long - node.long
	
			// Normalize the direction
			magnitude := math.Sqrt(latDirection*latDirection + longDirection*longDirection)
			latDirection /= magnitude
			longDirection /= magnitude
	
			// Move the node towards the hub
			node.lat += latDirection * distance
			node.long += longDirection * distance
		}
		
	} else {
		// Randomly choose a direction: 0 = up, 1 = down, 2 = left, 3 = right
		direction := rand.Intn(4)

		switch direction {
		case 0: // Move up (North)
			node.lat += distance
			//fmt.Println("Node moved up")
		case 1: // Move down (South)
			node.lat -= distance
			//fmt.Println("Node moved down")
		case 2: // Move left (West)
			node.long -= distance
			//fmt.Println("Node moved left")
		case 3: // Move right (East)
			node.long += distance
			//fmt.Println("Node moved right")
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./node host:port")
		return
	}

	// node hostname
	nodeHostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	plannerAddr := os.Args[1]

	listenPort := "8080"
	if len(os.Args) > 2 {
		listenPort = os.Args[2]
	}

	// Connect to Planner server
	conn, err := net.Dial("tcp", plannerAddr)
	if err != nil {
		fmt.Println("Error connecting to planner server:", err)
		return
	}

	msgHandler := messages.NewMessageHandler(conn)
	defer msgHandler.Close()
	defer conn.Close()

	// // pid is the nodeId
	nodeId := strconv.Itoa(os.Getpid())
	nodeAddr := nodeHostname + ":" + listenPort

	// Node coordinates
	lat, long := generateRandomCoordinate()

	// Battery level
	battery_level := generateBatteryLevel()
	seeking_charge := false
	if battery_level <= 30 {
		seeking_charge = true
	}

	node := nodeStruct{
		plannerAddr: plannerAddr,
		hubAddr: "",
		nodeAddr: nodeAddr,
		nodeId: nodeId,

		lat: lat,
		long: long,

		battery_level: battery_level,
		seeking_charge: seeking_charge,
	}

	// // send recon message to Planner
	_, reconHubAddr := sendRecon(&node)
	node.hubAddr = reconHubAddr

	// Connect to hub server
	conn, err = net.Dial("tcp", node.hubAddr)
	if err != nil {
		fmt.Println("Error connecting to hub server:", err)
		return
	}	

	// send init message
	_, nodeId = sendJoin(&node)

	hbChan := make(chan struct{})

	movePoint(&node)

	go listen(&node, hbChan)

	log.Printf("Node: %v sending heartbeats\n", node.nodeId)

	for {
		time.Sleep(5 * time.Second)
		err = sendHearbeats(&node, hbChan)
		if err != nil {
			return
		}
	}

}