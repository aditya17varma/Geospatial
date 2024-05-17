package node

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"

	"geospatial/messages"
)

const BATTERYCOST = 3

type NodeStruct struct {
	Mutex       sync.Mutex
	plannerAddr string
	hubAddr     string
	NodeAddr    string
	NodeId      string

	lat  float64
	long float64

	battery_level  int32
	seeking_charge bool
	charging       bool

	hub_set  bool
	hub_lat  float64
	hub_long float64
	at_hub   bool

	targetLat  float64
	targetLong float64
	at_target  bool
}

func sendRecon(node *NodeStruct) (bool, string) {
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

	currentTime := time.Now()
	// convert time to protobuf timestamp
	ts, _ := ptypes.TimestampProto(currentTime)

	reconMsg := messages.ReconRequest{NodeAddr: node.NodeAddr, NodeId: node.NodeId, Coordinates: &coordMsg, BatteryState: &batteryMsg, Timestamp: ts}
	wrapper := &messages.Wrapper{
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

func sendJoinPlanner(node *NodeStruct) (bool, string) {
	// connect to Planner
	conn, err := net.Dial("tcp", node.plannerAddr)
	if err != nil {
		log.Printf("Could not connect to planner:", node.plannerAddr)
		return false, "false"
	}
	defer conn.Close()

	msgHandler := messages.NewMessageHandler(conn)
	defer msgHandler.Close()

	// Send Join Request
	coordMsg := messages.Geocoordinate{Latitude: node.lat, Longitude: node.long}
	batteryMsg := messages.Battery{BatteryLevel: node.battery_level, SeekingCharge: node.seeking_charge}

	joinMsg := messages.JoinRequest{NodeAddr: node.NodeAddr, NodeId: node.NodeId, Coordinates: &coordMsg, BatteryState: &batteryMsg}
	wrapper := &messages.Wrapper{
		Msg: &messages.Wrapper_JoinReq{JoinReq: &joinMsg},
	}
	msgHandler.Send(wrapper)

	// Wait for Join Response
	wrapper, _ = msgHandler.Receive()
	msg := wrapper.GetJoinRes()
	if msg.Accept {
		log.Println("Node successfuly joined planner:", msg.NodeId)
		// hubCoordinates := msg.Coordinates
		// node.hub_set = true
		// node.hub_lat = hubCoordinates.Latitude
		// node.hub_long = hubCoordinates.Longitude

		return true, msg.NodeId
	} else {
		log.Println("Node could not join planner")
		return false, "false"
	}

	return false, "false"

}

// Sends join message to hub
func sendJoin(node *NodeStruct) (bool, string) {
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

	joinMsg := messages.JoinRequest{NodeAddr: node.NodeAddr, NodeId: node.NodeId, Coordinates: &coordMsg, BatteryState: &batteryMsg}
	wrapper := &messages.Wrapper{
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

func listen(node *NodeStruct, hbChan chan struct{}) {
	// listen for messages
	listener, err := net.Listen("tcp", node.NodeAddr)
	if err != nil {
		log.Printf("Could not listen on %s\n", node.NodeAddr)
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

func handleMessages(conn net.Conn, msgHandler *messages.MessageHandler, hbChan chan struct{}, node *NodeStruct) {
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
		case *messages.Wrapper_ServerMessage:
			serverMsg := msg.ServerMessage
			if serverMsg.Response == "ChargeComplete" {
				node.charging = false
				node.seeking_charge = false
				node.battery_level = 100
				node.at_hub = false
				node.at_target = false
				node.hub_set = false
				log.Println("Node charge complete")

				// Generate new target
				node.targetLat, node.targetLong = generateTargetCoordinate(node.hub_lat, node.hub_long)
			}
		default:
			log.Printf("Unexpected message type: %T", msg)
			return
		}
	}

}

func sendHearbeats(node *NodeStruct, hbChan <-chan struct{}) error {
	for {
		select {
		case <-hbChan:
			fmt.Println("Node killed, stopping heartbeats")
			return errors.New("node killed")
		default:

			// UPDATE NODE STATE (Coords & Battery)
			currentTime := time.Now()
			// convert time to protobuf timestamp
			ts, err := ptypes.TimestampProto(currentTime)
			if err != nil {
				fmt.Println("Error converting time to timestamp: %v\n", err)
				return err
			}

			node.Mutex.Lock()

			// move node
			movePoint(node)

			node.Mutex.Unlock()

			// Create heartbeat message
			coordMsg := messages.Geocoordinate{Latitude: node.lat, Longitude: node.long}
			batteryMsg := messages.Battery{BatteryLevel: node.battery_level, SeekingCharge: node.seeking_charge}
			hrbtMsg := messages.Heartbeat{NodeId: node.NodeId, State: "Alive", Timestamp: ts, Coordinates: &coordMsg, BatteryState: &batteryMsg}
			wrapper := &messages.Wrapper{
				Msg: &messages.Wrapper_HeartbeatMessage{HeartbeatMessage: &hrbtMsg},
			}

			if node.hub_set {
				// Send Heartbeat to Hub
				conn, err := net.Dial("tcp", node.hubAddr)
				if err != nil {
					log.Printf("Could not connect to hub: %s\n", node.hubAddr)
				}

				msgHandler := messages.NewMessageHandler(conn)
				defer msgHandler.Close()
				defer conn.Close()

				if err := msgHandler.Send(wrapper); err != nil {
					fmt.Println("Error sending hub heartbeat message:", err)
					return err
				}
			}

			// Send Heartbeat to Planner
			connP, err := net.Dial("tcp", node.plannerAddr)
			if err != nil {
				log.Printf("Could not connect to planner: %s\n", node.plannerAddr)
			}
			msgHandlerP := messages.NewMessageHandler(connP)
			defer msgHandlerP.Close()
			defer connP.Close()

			if err := msgHandlerP.Send(wrapper); err != nil {
				fmt.Println("Error sending planner heartbeat message:", err)
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

func generateTargetCoordinate(hubLat, hubLong float64) (float64, float64) {
	// Generate coordinates within 1km of the hub
	rand.Seed(time.Now().UnixNano())

	// Generate random adjustments
	deltaLat := (rand.Float64() - 0.5) * (0.01 * 2)
	deltaLong := (rand.Float64() - 0.5) * (0.01 * 2)

	// Apply randomness
	lat := hubLat + deltaLat
	long := hubLong + deltaLong

	return lat, long

}

// Returns a random number between 85 - 100 for starting battery level
func generateBatteryLevel() int32 {
	rand.Seed(time.Now().UnixNano())

	// Generate random number between 85 and 100
	randomNumber := rand.Intn(16) + 85

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
func movePoint(node *NodeStruct) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Define the small distanceStep to move in terms of degrees
	// This is a very rough approximation for demonstration purposes
	distanceStep := 0.001 // Roughly 111 meters for latitude, varies for longitude

	// If hub not set, join Hub
	if node.at_target && !node.hub_set {
		fmt.Printf("[Node %s] Hub not set, reconning Planner for Hub\n", node.NodeId)
		// send recon message to Planner
		_, reconHubAddr := sendRecon(node)
		node.hubAddr = reconHubAddr

		// Connect to hub server
		_, err := net.Dial("tcp", node.hubAddr)
		if err != nil {
			fmt.Println("Error connecting to hub server:", err)
			return
		}

		// send init message
		_, nodeId := sendJoin(node)
		fmt.Printf("[Node %s] Node Joined\n", nodeId)
	}

	// fmt.Println("MOVE POINT")
	// fmt.Printf("[Node %s] charging: %v sc: %v at_hub: %v\n", node.NodeId, node.charging, node.seeking_charge, node.at_hub)
	if !node.charging { // Moving towards hub or target
		if node.seeking_charge && !node.at_hub { // Move towards hub

			// fmt.Println("Moving towards hub")
			hubDistance := haversine(node.lat, node.long, node.hub_lat, node.hub_long)
			hubDistanceMeters := hubDistance * 1000

			if hubDistanceMeters <= 100 { // Finer movement when close to hub
				distanceStep = 0.0001 // Prevent overshooting and oscillation
			}

			if hubDistanceMeters <= 30 { // At hub don't move
				log.Printf("[Node %s] Node at hub, not moving\n", node.NodeId)
				node.at_hub = true
				node.charging = true
				node.seeking_charge = false

			} else {
				fmt.Printf("[Node %s] Moving towards hub, distance: %f battery: %d\n", node.NodeId, hubDistanceMeters, node.battery_level)
				// update battery
				node.battery_level = node.battery_level - BATTERYCOST

				// Calculate the direction to move in
				latDirection := node.hub_lat - node.lat
				longDirection := node.hub_long - node.long

				// Normalize the direction
				magnitude := math.Sqrt(latDirection*latDirection + longDirection*longDirection)
				latDirection /= magnitude
				longDirection /= magnitude

				// Move the node towards the hub
				node.lat += latDirection * distanceStep
				node.long += longDirection * distanceStep
			}

		} else { // Move towards target
			// fmt.Println("Moving towards target")
			// Distance to target
			targetDistance := haversine(node.lat, node.long, node.targetLat, node.targetLong)
			targetDistanceMeters := targetDistance * 1000

			if targetDistanceMeters <= 100 { // Finer movement when close to target
				distanceStep = 0.0001 // Prevent overshooting and oscillation
			}

			if targetDistanceMeters <= 30 { // At target
				node.at_target = true
				fmt.Printf("[Node %s] Node at target\n", node.NodeId)
				node.seeking_charge = true

				// Message Planner for Hub assignment
				conn, err := net.Dial("tcp", node.plannerAddr)
				if err != nil {
					log.Printf("Could not connect to planner:", node.plannerAddr)
					return
				}
				defer conn.Close()

			} else { // Move towards target
				// update battery
				node.battery_level = node.battery_level - BATTERYCOST

				fmt.Printf("[Node %s] Moving towards target, distance: %f battery: %d\n", node.NodeId, targetDistanceMeters, node.battery_level)
				// Calculate the direction to move in
				latDirection := node.targetLat - node.lat
				longDirection := node.targetLong - node.long

				// Normalize the direction
				magnitude := math.Sqrt(latDirection*latDirection + longDirection*longDirection)
				latDirection /= magnitude
				longDirection /= magnitude

				// Move the node towards the target
				node.lat += latDirection * distanceStep
				node.long += longDirection * distanceStep
			}
		}
	} else {
		fmt.Printf("[Node %s] Node charging\n", node.NodeId)
	}

}

// Setup Node
func NewNode(plannerAddress, nodePort string, startingLat, stratingLong float64) (*NodeStruct, error) {
	// node hostname
	nodeHostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	plannerAddr := plannerAddress

	listenPort := nodePort

	// // pid is the nodeId
	nodeId := strconv.Itoa(os.Getpid())
	nodeAddr := nodeHostname + ":" + listenPort

	targetLat, targetLong := generateTargetCoordinate(startingLat, stratingLong)

	node := NodeStruct{
		plannerAddr: plannerAddr,
		hubAddr:     "",
		NodeAddr:    nodeAddr,
		NodeId:      nodeId,

		lat:  startingLat,
		long: stratingLong,

		battery_level:  100,
		seeking_charge: false,
		charging:       false,

		targetLat:  targetLat,
		targetLong: targetLong,

		at_target: false,
		at_hub:    false,
		hub_set:   false,
	}

	return &node, nil
}

// Function to activate a node in geospatial/main, instead of as a sepearate process
// For simulation purposes
func (node *NodeStruct) ActivateNode() {
	fmt.Printf("[Node %s] Activating Node %s\n", node.NodeId, node.NodeAddr)

	// listenPort := nodePort

	// Connect to Planner server
	conn, err := net.Dial("tcp", node.plannerAddr)
	if err != nil {
		fmt.Println("Error connecting to planner server:", err)
		return
	}

	msgHandler := messages.NewMessageHandler(conn)
	defer msgHandler.Close()
	defer conn.Close()

	// Register with Planner
	sendJoinPlanner(node)

	hbChan := make(chan struct{})

	movePoint(node)

	go listen(node, hbChan)

	log.Printf("[Node %s] Node sending heartbeats\n", node.NodeId)

	for {
		time.Sleep(2 * time.Second)
		err = sendHearbeats(node, hbChan)
		if err != nil {
			return
		}
	}
}

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: ./node host:port")
// 		return
// 	}

// 	// node hostname
// 	nodeHostname, err := os.Hostname()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	plannerAddr := os.Args[1]

// 	listenPort := "8080"
// 	if len(os.Args) > 2 {
// 		listenPort = os.Args[2]
// 	}

// 	// Connect to Planner server
// 	conn, err := net.Dial("tcp", plannerAddr)
// 	if err != nil {
// 		fmt.Println("Error connecting to planner server:", err)
// 		return
// 	}

// 	msgHandler := messages.NewMessageHandler(conn)
// 	defer msgHandler.Close()
// 	defer conn.Close()

// 	// // pid is the nodeId
// 	nodeId := strconv.Itoa(os.Getpid())
// 	nodeAddr := nodeHostname + ":" + listenPort

// 	// Node coordinates
// 	lat, long := generateRandomCoordinate()

// 	// Battery level
// 	battery_level := generateBatteryLevel()
// 	seeking_charge := false
// 	if battery_level <= 30 {
// 		seeking_charge = true
// 	}

// 	node := nodeStruct{
// 		plannerAddr: plannerAddr,
// 		hubAddr:     "",
// 		nodeAddr:    nodeAddr,
// 		nodeId:      nodeId,

// 		lat:  lat,
// 		long: long,

// 		battery_level:  battery_level,
// 		seeking_charge: seeking_charge,
// 	}

// 	// // send recon message to Planner
// 	_, reconHubAddr := sendRecon(&node)
// 	node.hubAddr = reconHubAddr

// 	// Connect to hub server
// 	conn, err = net.Dial("tcp", node.hubAddr)
// 	if err != nil {
// 		fmt.Println("Error connecting to hub server:", err)
// 		return
// 	}

// 	// send init message
// 	_, nodeId = sendJoin(&node)

// 	hbChan := make(chan struct{})

// 	movePoint(&node)

// 	go listen(&node, hbChan)

// 	log.Printf("Node: %v sending heartbeats\n", node.nodeId)

// 	for {
// 		time.Sleep(5 * time.Second)
// 		err = sendHearbeats(&node, hbChan)
// 		if err != nil {
// 			return
// 		}
// 	}

// }
