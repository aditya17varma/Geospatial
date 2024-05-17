package hub

// protoc --go_out=. proto/heartbeat.proto

import (
	"fmt"
	"net"
	//"os"
	"container/heap"
	"log"
	"math"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/paulmach/orb"

	"geospatial/messages"
)

type Node struct {
	Id            string
	State         string
	Timestamp     *timestamp.Timestamp
	NodeAddr      string
	Coordinates   Coordinate
	Bucket        int
	BatteryLevel  int32
	SeekingCharge bool
	AtHub         bool
	Distance      float64
	Charging      bool
	DowntimeStart time.Time
	HubAddr       string
	ChargingSpot  int
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Hub struct {
	Mutex   sync.Mutex
	NodeMap map[string]Node
	idCount int64

	Coordinates Coordinate

	// nodeBucket1 map[string]Node
	// nodeBucket2 map[string]Node
	// nodeBucket3 map[string]Node

	nodeQueue PriorityQueue

	HubPoint orb.Point

	HubHost string
	HubPort string

	AllocationPolicy string

	HubChargingSpotsTotal     int
	HubChargingSpotsAvailable int
	NodesChargedTotal         int

	HubNodesCharging map[string]Node
	HubNodesWaiting  map[string]Node

	HubDowntime time.Duration
}

// PriorityQueue implements heap.Interface and holds Nodes.
type PriorityQueue struct {
	nodes            []*Node
	allocationPolicy string
}

func (pq PriorityQueue) Len() int { return len(pq.nodes) }

func (pq PriorityQueue) Less(i, j int) bool {
	switch pq.allocationPolicy {
	case "minBattery":
		return pq.nodes[i].BatteryLevel < pq.nodes[j].BatteryLevel
	case "minDistance":
		return pq.nodes[i].Distance < pq.nodes[j].Distance
	default:
		return pq.nodes[i].Distance < pq.nodes[j].Distance
	}
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.nodes[i], pq.nodes[j] = pq.nodes[j], pq.nodes[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	pq.nodes = append(pq.nodes, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.nodes
	n := len(old)
	item := old[n-1]
	pq.nodes = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) updatePolicy(policy string) {
	pq.allocationPolicy = policy
	heap.Init(pq) // re-heapify the queue based on the new policy
}

func NewHub(hubHost, hubPort, allocationPolicy string, lat, long float64, chargingSpots int) *Hub {
	coords := Coordinate{
		Latitude:  lat,
		Longitude: long,
	}

	hp := orb.Point{lat, long}

	// fmt.Println("Init hub with charging spots: ", chargingSpots)

	return &Hub{
		NodeMap:     make(map[string]Node),
		idCount:     0,
		Coordinates: coords,

		HubPoint: hp,

		HubHost: hubHost,
		HubPort: hubPort,

		AllocationPolicy: allocationPolicy,

		HubChargingSpotsTotal:     chargingSpots,
		HubChargingSpotsAvailable: chargingSpots,
		NodesChargedTotal:         0,

		HubNodesCharging: make(map[string]Node),
		HubNodesWaiting:  make(map[string]Node),

		nodeQueue: PriorityQueue{
			nodes:            make([]*Node, 0),
			allocationPolicy: allocationPolicy,
		},
	}
}

// Report the state of the hub
// Number of charging spots available, number of nodes waiting, number of nodes allocated
func (hub *Hub) ReportHubState() []int {
	sizes := make([]int, 3)

	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	sizes[0] = hub.HubChargingSpotsAvailable
	sizes[1] = len(hub.HubNodesWaiting)
	sizes[2] = len(hub.NodeMap)

	return sizes
}

func (hub *Hub) addNode(node Node) string {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	// Add node to hub's NodeMap
	hub.NodeMap[node.Id] = node

	return node.Id
}

func (hub *Hub) updateState(nodeId, state string) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	node.State = state
	hub.NodeMap[nodeId] = node
}

func (hub *Hub) updateTimestamp(nodeId string, timestamp *timestamp.Timestamp, lat, long float64, battery_level int32, seeking_charge bool) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	node.Timestamp = timestamp
	node.Coordinates.Latitude = lat
	node.Coordinates.Longitude = long
	node.BatteryLevel = battery_level
	node.SeekingCharge = seeking_charge
	// hub.NodeMap[nodeId] = node

	distance := hub.haversine(lat, long)
	nodeDistanceMeters := distance * 1000
	//log.Println("Node distance:", distance)
	node.Distance = nodeDistanceMeters

	// Check if Node at Hub
	if nodeDistanceMeters <= 30 {
		node.AtHub = true
		fmt.Println("[Hub] Node", nodeId, "at hub")
		// Check if Charging Spots Available
		if !node.Charging {
			if hub.HubChargingSpotsAvailable > 0 {
				// Charge Node
				node.Charging = true
				hub.HubChargingSpotsAvailable--
				hub.HubNodesCharging[nodeId] = node
				fmt.Println("[Hub] Charging available, charging node:", nodeId)
			} else {
				// Add Node to Waiting Queue, start downtime
				node.DowntimeStart = time.Now()
				heap.Push(&hub.nodeQueue, &node)
				hub.HubNodesWaiting[nodeId] = node
				fmt.Println("[Hub] No charging available, adding node to waiting queue:", nodeId)
			}
		} else {
			fmt.Println("[Hub] Node", nodeId, "already charging")
		}

	}

	hub.NodeMap[nodeId] = node

}

func (hub *Hub) updateCoordinates(nodeId string, lat, long float64) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	node := hub.NodeMap[nodeId]
	node.Coordinates.Latitude = lat
	node.Coordinates.Longitude = long
	hub.NodeMap[nodeId] = node
}

func (hub *Hub) getCoordinates(nodeId string) (float64, float64) {
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

func killNode(nodeAddr string) {
	// connect to node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	msgHandler := messages.NewMessageHandler(conn)
	killMsg := messages.KillNode{Kill: true}
	wrapper := &messages.Wrapper{
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

			// fmt.Println("Hearbeat from:", node_id)
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

func (hub *Hub) handleJoin(msgHandler *messages.MessageHandler, nodeAddr, nodeId string, lat, long float64, battery_level int32, seeking_charge bool) {
	currentTime := time.Now()

	// convert time to protobuf timestamp
	ts, err := ptypes.TimestampProto(currentTime)
	if err != nil {
		fmt.Println("Error converting time to timestamp: %v\n", err)
		return
	}

	coordinates := Coordinate{
		Latitude:  lat,
		Longitude: long,
	}

	// distance in km
	distance := hub.haversine(lat, long)

	newNode := Node{
		State:         "Join",
		Timestamp:     ts,
		NodeAddr:      nodeAddr,
		Id:            nodeId,
		Coordinates:   coordinates,
		Bucket:        0,
		BatteryLevel:  battery_level,
		SeekingCharge: seeking_charge,
		AtHub:         false,
		Distance:      distance,
		Charging:      false,
	}
	// add node to map
	nodeId = hub.addNode(newNode)

	hubCoordMsg := messages.Geocoordinate{Latitude: hub.Coordinates.Latitude, Longitude: hub.Coordinates.Longitude}

	joinRes := messages.JoinResponse{Accept: true, NodeId: nodeId, Coordinates: &hubCoordMsg}

	wrapper := &messages.Wrapper{
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
		wrapper := &messages.Wrapper{
			Msg: &messages.Wrapper_ServerMessage{ServerMessage: &serverMsg},
		}
		msgHandler.Send(wrapper)
		fmt.Println("Server sent Join to node")
	} else { //node exists
		// check state

		if tempNode.State == "Dead" {
			// send reinit
			fmt.Println("Dead node, sending kill")
			killNode(tempNode.NodeAddr)
		} else if tempNode.State == "Alive" {
			// check time for timeout
			prevTime := tempNode.Timestamp
			//timeDiff := nodeTimestamp.Seconds - prevTime.Seconds
			if (nodeTimestamp.Seconds - prevTime.Seconds) > 10 {
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
		} else if tempNode.State == "Join" {
			// update nodeMap
			hub.updateState(nodeId, "Alive")
			hub.updateTimestamp(nodeId, nodeTimestamp, lat, long, battery_level, seeking_charge)

			//fmt.Println("Node state from Join to Alive")
		}
	}
}

// Check Hub charging spots periodically and move nodes from queue to charging
func (hub *Hub) checkHubState() {
	for {
		time.Sleep(5 * time.Second)

		hub.Mutex.Lock()

		for nodeId, node := range hub.HubNodesCharging {
			node.BatteryLevel += 10
			// Check if Node is Fully Charged
			if node.BatteryLevel >= 85 {
				fmt.Println("[Hub] Node", nodeId, "fully charged")
				node.Charging = false

				hub.HubChargingSpotsAvailable++
				hub.NodesChargedTotal += 1


				delete(hub.HubNodesCharging, nodeId)
				delete(hub.NodeMap, nodeId)
				sendChargeUpdate(node.NodeAddr)

				// Check if Nodes waiting for Charging
				if hub.nodeQueue.Len() > 0 {
					nextNode := heap.Pop(&hub.nodeQueue).(*Node)
					nextNode.Charging = true
					hub.HubNodesCharging[nextNode.Id] = *nextNode
					hub.HubChargingSpotsAvailable--
					delete(hub.HubNodesWaiting, nextNode.Id)
					fmt.Println("[Hub] Node", nextNode.Id, "moved from queue to charging")
					// Add downtime
					nodeDowntime := time.Since(nextNode.DowntimeStart)
					hub.HubDowntime += nodeDowntime
				}
			} else {
				fmt.Println("[Hub] Node", nodeId, "charging, battery level:", node.BatteryLevel)
				hub.HubNodesCharging[nodeId] = node
			}
		}
		hub.Mutex.Unlock()
	}
}

// Send message to node that it is done charging
func sendChargeUpdate(nodeAddr string) {
	// connect to node
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	msgHandler := messages.NewMessageHandler(conn)
	serverMsg := messages.ServerMessage{Response: "ChargeComplete"}
	wrapper := &messages.Wrapper{
		Msg: &messages.Wrapper_ServerMessage{ServerMessage: &serverMsg},
	}
	msgHandler.Send(wrapper)
	log.Println("[Hub] Sent charge update to:", nodeAddr)
	return
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

	// Start Hub State Checker
	go hub.checkHubState()

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

func (hub *Hub) GetAverageDowntimeStat() time.Duration {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	if hub.NodesChargedTotal == 0 {
		return 0
	}

	return hub.HubDowntime / time.Duration(hub.NodesChargedTotal)
}

func (hub *Hub) GetHubDowntime() time.Duration {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	return hub.HubDowntime
}
