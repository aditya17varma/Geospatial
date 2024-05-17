package planner

import (
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"

	"geospatial/hub"
	"geospatial/messages"
)

type Planner struct {
	Mutex sync.Mutex

	HubMap           map[string]*hub.Hub
	NodeMap          map[string]*hub.Node
	AllocationPolicy string
}

// Return new Planner struct
// Default AllocationPolicy is "distance"
func NewPlanner(allocationPolicy string) *Planner {
	return &Planner{
		HubMap:           make(map[string]*hub.Hub),
		NodeMap:          make(map[string]*hub.Node),
		AllocationPolicy: allocationPolicy,
	}
}

func (planner *Planner) RegisterHub(hub *hub.Hub) {
	planner.Mutex.Lock()
	hubName := hub.HubHost + ":" + hub.HubPort
	planner.HubMap[hubName] = hub
	planner.Mutex.Unlock()

	log.Println("Registered Hub: %s", hubName)
}

func (planner *Planner) SetAllocationPolicy(policy string) {
	planner.AllocationPolicy = policy
}

func (planner *Planner) HandleMessage(conn net.Conn, msgHandler *messages.MessageHandler) {
	defer msgHandler.Close()
	defer conn.Close()

	for {
		wrapper, _ := msgHandler.Receive()

		switch msg := wrapper.Msg.(type) {
		case *messages.Wrapper_JoinReq:
			// connect to hub
			initMsg := msg.JoinReq

			nodeAddr := initMsg.NodeAddr
			nodeId := initMsg.NodeId
			coordMsg := initMsg.Coordinates
			lat := coordMsg.Latitude
			long := coordMsg.Longitude
			batteryMsg := initMsg.BatteryState
			batteryLevel := batteryMsg.BatteryLevel
			seekingCharge := batteryMsg.SeekingCharge

			// convert time to protobuf timestamp
			ts, err := ptypes.TimestampProto(time.Now())
			if err != nil {
				fmt.Println("Error converting time to timestamp: %v\n", err)
				return
			}

			coordinates := hub.Coordinate{
				Latitude:  lat,
				Longitude: long,
			}

			// Add node to NodeMap
			planner.Mutex.Lock()
			newNode := hub.Node{
				State:         "Join",
				Timestamp:     ts,
				NodeAddr:      nodeAddr,
				Id:            nodeId,
				Coordinates:   coordinates,
				Bucket:        0,
				BatteryLevel:  batteryLevel,
				SeekingCharge: seekingCharge,
				AtHub:         false,
				Distance:      0,
				Charging:      false,
			}
			planner.NodeMap[nodeId] = &newNode

			planner.Mutex.Unlock()

			// Send response to node
			hubCoordMsg := messages.Geocoordinate{Latitude: coordinates.Latitude, Longitude: coordinates.Longitude}

			joinRes := messages.JoinResponse{Accept: true, NodeId: nodeId, Coordinates: &hubCoordMsg}
			wrapper := &messages.Wrapper{
				Msg: &messages.Wrapper_JoinRes{JoinRes: &joinRes},
			}
			msgHandler.Send(wrapper)

		case *messages.Wrapper_ReconReq:
			reconMsg := msg.ReconReq

			node_addr := reconMsg.NodeAddr
			node_id := reconMsg.NodeId

			coordMsg := reconMsg.Coordinates
			lat := coordMsg.Latitude
			long := coordMsg.Longitude

			batteryMsg := reconMsg.BatteryState
			battery_level := batteryMsg.BatteryLevel
			seeking_charge := batteryMsg.SeekingCharge

			planner.assignHub(msgHandler, node_addr, node_id, lat, long, battery_level, seeking_charge)

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

			// fmt.Println("Planner received heartbeat from node:", node_id)

			planner.Mutex.Lock()

			node := planner.NodeMap[node_id]
			node.State = state
			node.Timestamp = timestamp
			node.Coordinates.Latitude = lat
			node.Coordinates.Longitude = long
			node.BatteryLevel = battery_level
			node.SeekingCharge = seeking_charge

			planner.Mutex.Unlock()

		case nil:
			return

		default:
			fmt.Printf("Incorrect message type from Client %T\n", msg)
			return
		}
	}
}

func (planner *Planner) assignHub(msgHandler *messages.MessageHandler, node_addr, node_id string, lat, long float64, battery_level int32, seeking_charge bool) {
	//log.Println("Assigning hub to node:", node_id)

	planner.Mutex.Lock()
	defer planner.Mutex.Unlock()

	var assignedHub *hub.Hub
	var assignedKey string

	if planner.AllocationPolicy == "minDistance" { // Find min distance hub
		fmt.Println("Using minDistance policy")
		minDistance := math.MaxFloat64
		var closestHub *hub.Hub
		var closestKey string

		for key, hub := range planner.HubMap {
			distance := haversine(lat, long, hub.Coordinates.Latitude, hub.Coordinates.Longitude)
			if distance < minDistance {
				closestHub = hub
				minDistance = distance
				closestKey = key
			}
		}
		assignedHub = closestHub
		assignedKey = closestKey
	}

	// Update hubAddr
	node := planner.NodeMap[node_id]
	node.HubAddr = assignedHub.HubHost + ":" + assignedHub.HubPort
	planner.NodeMap[node_id] = node

	log.Printf("Assigning node %+v to hub: %+v", node_id, assignedKey)

	reconRes := messages.ReconResponse{HubHost: assignedHub.HubHost, HubPort: assignedHub.HubPort, NodeId: node_id, Accept: true}

	wrapper := &messages.Wrapper{
		Msg: &messages.Wrapper_ReconRes{ReconRes: &reconRes},
	}
	msgHandler.Send(wrapper)

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

// Start the Planner server
func (planner *Planner) StartListening(listenPort string) {
	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		log.Fatalf("Error listening: %s", err)
	}
	defer listener.Close()

	log.Println("Planner server started on port:", listenPort)

	for {
		// Accept incoming connections to Planner
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Message Handler
		msgHandler := messages.NewMessageHandler(conn)
		go planner.HandleMessage(conn, msgHandler)
	}
}
