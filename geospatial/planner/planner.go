package planner

import (
	"sync"
	"fmt"
	"math"
	"log"
	"net"
	
	"geospatial/messages"
	"geospatial/hub"
)


type Planner struct {
	Mutex sync.Mutex
	
	HubMap map[string]*hub.Hub
	AllocationPolicy string	
}

// Return new Planner struct
// Default AllocationPolicy is "distance"
func NewPlanner() *Planner {
	return &Planner {
		HubMap: make(map[string]*hub.Hub),
		AllocationPolicy: "distance",
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

func (planner *Planner) HandleMessage (conn net.Conn, msgHandler *messages.MessageHandler) {
	defer msgHandler.Close()
	defer conn.Close()
	
	for {
		wrapper, _ := msgHandler.Receive()

		switch msg := wrapper.Msg.(type) {
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

	// Find min distance hub
	minDistance := math.MaxFloat64
	var closestHub *hub.Hub;
	var closestKey string;

	for key, hub := range planner.HubMap {
		distance := haversine(lat, long, hub.Coordinates.Latitude, hub.Coordinates.Longitude)
		if distance < minDistance {
			closestHub = hub
			minDistance = distance
			closestKey = key
		}
	}

	log.Printf("Assigning node %+v to hub: %+v", node_id, closestKey)

	reconRes := messages.ReconResponse{HubHost: closestHub.HubHost, HubPort: closestHub.HubPort, NodeId: node_id, Accept: true}
	
	wrapper := &messages.Wrapper {
		Msg: &messages.Wrapper_ReconRes{ReconRes: &reconRes},
	}
	msgHandler.Send(wrapper)

	return
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