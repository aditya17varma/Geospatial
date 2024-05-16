package hub

// protoc --go_out=. proto/heartbeat.proto 

import (
	"geospatial/messages"
	"fmt"
	"net"
	//"os"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes"
	"sync"
	"time"
	"log"
	"math"

	"github.com/paulmach/orb"
)

type Node struct {
	Id string
	State string
	Timestamp *timestamp.Timestamp
	NodeAddr string
	Coordinates Coordinate
	Bucket int
	BatteryLevel int32
	SeekingCharge bool
}

type Coordinate struct {
	Latitude float64
	Longitude float64
}

type Hub struct {
	Mutex sync.Mutex
	NodeMap map[string]Node
	idCount int64

	Coordinates Coordinate

	nodeBucket1 map[string]Node
	nodeBucket2 map[string]Node
	nodeBucket3 map[string]Node

	HubPoint orb.Point

	HubHost string
	HubPort string
}

func NewHub(hubHost, hubPort string, lat, long float64) *Hub {
	coords := Coordinate{
		Latitude: lat,
		Longitude: long,
	}

	hp := orb.Point{lat, long}

	return &Hub {
		NodeMap: make(map[string]Node),
		idCount : 0,
		Coordinates: coords,

		nodeBucket1: make(map[string]Node),
		nodeBucket2: make(map[string]Node),
		nodeBucket3: make(map[string]Node),

		HubPoint: hp,

		HubHost: hubHost,
		HubPort: hubPort,
	}
}

func (hub *Hub) ReportBucketSizes() []int {
	sizes := make([]int, 3)

	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	sizes[0] = len(hub.nodeBucket1)
	sizes[1] = len(hub.nodeBucket2)
	sizes[2] = len(hub.nodeBucket3)

	return sizes
}

func (hub *Hub) addNode(node Node) string {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	lat := node.Coordinates.Latitude
	long := node.Coordinates.Longitude

	// distance in km
	distance := hub.haversine(lat, long)
	//log.Println("Distance:", distance)

	if distance <= 1 {
		hub.nodeBucket1[node.Id] = node
		node.Bucket = 1
	} else if distance <= 2 {
		hub.nodeBucket2[node.Id] = node
		node.Bucket = 2
	} else if distance <= 3 {
		hub.nodeBucket3[node.Id] = node
		node.Bucket = 3
	}

	hub.NodeMap[node.Id] = node

	//log.Println("Node added to bucket:", node.Bucket)

	return node.Id
}

func (hub *Hub) updateState(nodeId, state string) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	node.State = state
	hub.NodeMap[nodeId] = node
}

func (hub *Hub) updateTimestamp(nodeId string, timestamp *timestamp.Timestamp, lat, long float64, battery_level int32, seeking_charge bool){
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	node.Timestamp = timestamp
	node.Coordinates.Latitude = lat
	node.Coordinates.Longitude = long
	node.BatteryLevel = battery_level
	node.SeekingCharge = seeking_charge
	hub.NodeMap[nodeId] = node

	distance := hub.haversine(lat, long)
	//log.Println("Node distance:", distance)
	currentBucket := node.Bucket

	//log.Println("CurrentBucket:", currentBucket)

	// Bucket Distances: 1km, 2km, 3km
	if distance <= 1 {
		if currentBucket == 0 {
			// Not assigned to any bucket
			node.Bucket = 1
		} else if currentBucket != 1 {
			// remove from other bucket
			if currentBucket == 2 {
				delete(hub.nodeBucket2, node.Id)
			} else {
				delete(hub.nodeBucket3, node.Id)
			}
		}
		hub.nodeBucket1[node.Id] = node
	} else if distance <= 2 {
		if currentBucket == 0 {
			// Not assigned to any bucket
			node.Bucket = 2
		} else if currentBucket != 2 {
			// remove from other bucket
			if currentBucket == 1 {
				delete(hub.nodeBucket1, node.Id)
			} else {
				delete(hub.nodeBucket3, node.Id)
			}
		}
		hub.nodeBucket2[node.Id] = node
	} else if distance <= 3 {
		if currentBucket == 0 {
			// Not assigned to any bucket
			node.Bucket = 3
		} else if currentBucket != 3 {
			// remove from other bucket
			if currentBucket == 1 {
				delete(hub.nodeBucket1, node.Id)
			} else {
				delete(hub.nodeBucket2, node.Id)
			}
		}
		hub.nodeBucket3[node.Id] = node
	}
	hub.NodeMap[nodeId] = node

}

func (hub *Hub) updateCoordinates (nodeId string, lat, long float64) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	node.Coordinates.Latitude = lat
	node.Coordinates.Longitude = long
	hub.NodeMap[nodeId] = node
}

func (hub *Hub) getCoordinates (nodeId string) (float64, float64) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	lat := node.Coordinates.Latitude
	long := node.Coordinates.Longitude
	
	return lat, long

}

// haversine function calculates the distance between two points on the Earth
// given their latitude and longitude in degrees.
func (hub *Hub) haversine(nodeLat, nodeLong float64) float64 {
	// convert degrees to radians
	var dlat, dlon, a, c float64
	rad := math.Pi / 180
	
	hubLat := hub.Coordinates.Latitude
	hubLong := hub.Coordinates.Longitude

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

func killNode(nodeAddr string){
	// connect to node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	msgHandler := messages.NewMessageHandler(conn)
	killMsg := messages.KillNode{Kill: true}
	wrapper := &messages.Wrapper {
		Msg: &messages.Wrapper_KillNode{KillNode: &killMsg},
	}
	msgHandler.Send(wrapper)
	//log.Println("Sent kill to:", nodeAddr)
	return
}

func (hub *Hub) HandleMessage(conn net.Conn, msgHandler *messages.MessageHandler) {
	defer msgHandler.Close()
	defer conn.Close()
	
	for {
		wrapper, _ := msgHandler.Receive()

		switch msg := wrapper.Msg.(type) {
		case *messages.Wrapper_JoinReq:
			initMsg := msg.JoinReq

			node_addr := initMsg.NodeAddr
			nodeId := initMsg.NodeId

			coordMsg := initMsg.Coordinates
			lat := coordMsg.Latitude
			long := coordMsg.Longitude

			batteryMsg := initMsg.BatteryState
			battery_level := batteryMsg.BatteryLevel
			seeking_charge := batteryMsg.SeekingCharge

			//fmt.Println("Node coordinates: %d, %d\n", lat, long)

			hub.handleJoin(msgHandler, node_addr, nodeId, lat, long, battery_level, seeking_charge)

		case *messages.Wrapper_HeartbeatMessage:
			// Accessing the Heartbeat message from the wrapper
			
			hrbtMsg := msg.HeartbeatMessage

			// Accessing the node data in the Heartbeat message
			node_id := hrbtMsg.NodeId
			state := hrbtMsg.State
			timestamp := hrbtMsg.Timestamp
			
			coordMsg := hrbtMsg.Coordinates
			lat := coordMsg.Latitude
			long := coordMsg.Longitude

			batteryMsg := hrbtMsg.BatteryState
			battery_level := batteryMsg.BatteryLevel
			seeking_charge := batteryMsg.SeekingCharge

			fmt.Println("Hearbeat from:", node_id)
			//fmt.Println("BatterLvl:", battery_level)
			//fmt.Println("Seeking Charge:", seeking_charge)

			hub.handleHearbeat(msgHandler, node_id, state, timestamp, lat, long, battery_level, seeking_charge)
			
		case nil:
			return
		
		default:
			fmt.Printf("Incorrect message type from Client %T\n", msg)
			return
		}
	}
	
}

func (hub *Hub) handleJoin(msgHandler *messages.MessageHandler, nodeAddr, nodeId string, lat, long float64, battery_level int32, seeking_charge bool){
	currentTime := time.Now()

	// convert time to protobuf timestamp
	ts, err := ptypes.TimestampProto(currentTime)
	if err != nil {
		fmt.Println("Error converting time to timestamp: %v\n", err)
		return
	}

	coordinates := Coordinate{
		Latitude: lat,
		Longitude: long,
	}
	
	newNode := Node{
		State: "Join",
		Timestamp: ts,
		NodeAddr: nodeAddr,
		Id: nodeId,
		Coordinates: coordinates,
		Bucket: 0,
		BatteryLevel: battery_level,
		SeekingCharge: seeking_charge,
	}
	// add node to map
	nodeId = hub.addNode(newNode)

	hubCoordMsg := messages.Geocoordinate{Latitude: hub.Coordinates.Latitude, Longitude: hub.Coordinates.Longitude}

	joinRes := messages.JoinResponse{Accept: true, NodeId: nodeId, Coordinates: &hubCoordMsg}
	
	wrapper := &messages.Wrapper {
		Msg: &messages.Wrapper_JoinRes{JoinRes: &joinRes},
	}
	// fmt.Println("Sending:", wrapper.GetJoinRes())
	msgHandler.Send(wrapper)
	// fmt.Println("Server sent Join success")
	return

}

func (hub *Hub) handleHearbeat(msgHandler *messages.MessageHandler, nodeId, nodeState string, nodeTimestamp *timestamp.Timestamp, lat, long float64,
	battery_level int32, seeking_charge bool) {
	// check if nodeId in nodeMap map
	hub.Mutex.Lock()
	tempNode, exists := hub.NodeMap[nodeId]
	hub.Mutex.Unlock()
	if !exists { // node needs to initialize
		fmt.Println("Uninitialized node:", nodeId)
		serverMsg := messages.ServerMessage{Response: "Join"}
		wrapper := &messages.Wrapper {
			Msg: &messages.Wrapper_ServerMessage{ServerMessage: &serverMsg},
		}
		msgHandler.Send(wrapper)
		fmt.Println("Server sent Join to node")
	} else { //node exists
		// check state

		if (tempNode.State == "Dead"){
			// send reinit
			fmt.Println("Dead node, sending kill")
			killNode(tempNode.NodeAddr)
		} else if (tempNode.State == "Alive"){
			// check time for timeout
			prevTime := tempNode.Timestamp
			//timeDiff := nodeTimestamp.Seconds - prevTime.Seconds
			if ((nodeTimestamp.Seconds - prevTime.Seconds) > 10){
				// node timed out
				fmt.Println("Node timed out, kill")
				hub.updateState(nodeId, "Dead")
				hub.updateCoordinates(nodeId, lat, long)
				killNode(tempNode.NodeAddr)
				fmt.Println("Server sent kill to node")

			} else {
				hub.updateTimestamp(nodeId, nodeTimestamp, lat, long, battery_level, seeking_charge)
				// fmt.Println("Updated timestamp for:", nodeId)
				// fmt.Printf("Node %s moved to %f, %f\n", nodeId, lat, long)
			}
		} else if (tempNode.State == "Join") {
			// update nodeMap
			hub.updateState(nodeId, "Alive")
			hub.updateTimestamp(nodeId, nodeTimestamp, lat, long, battery_level, seeking_charge)

			//fmt.Println("Node state from Join to Alive")
		}
	}
}

func (hub *Hub) Listen() {
	hubAddr := hub.HubHost + ":" + hub.HubPort

	// Listen for incoming connections
	listener, err := net.Listen("tcp", ":"+hub.HubPort)
	if err != nil {
		log.Printf("Error starting hub server %s: %+v\n", hubAddr, err)
		return
	}

	log.Println("Hub listening:", hubAddr)

	defer listener.Close()

	for {
		// Accept incoming connections to Hub
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection %s: %+v\n", hubAddr, err)
			continue
		}
		// Message Handler
		msgHandler := messages.NewMessageHandler(conn)
		//go hub.HandleMessage(conn, msgHandler, hubStation)
		go hub.HandleMessage(conn, msgHandler)
	}

}
