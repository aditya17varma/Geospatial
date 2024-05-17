package customMap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/morikuni/go-geoplot"

	"geospatial/hub"
	"geospatial/planner"
)

func StartMapServer(customMap *CustomMap) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Convert CustomMap to geoplot.Map just before serving
		gMap := customMap.ToGeoPlotMap()
		err := geoplot.ServeMap(w, r, gMap)
		if err != nil {
			log.Println("Failed to serve map:", err)
			http.Error(w, "Failed to serve map", http.StatusInternalServerError)
		}
	})

	log.Println("Starting web server on :9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Println("Failed to start web server:", err)
	}
}

type CustomMap struct {
	Center *geoplot.LatLng
	Zoom   int
	Area   *geoplot.Area
	// Keep your own list of markers
	Markers []*geoplot.Marker
	Circles []*geoplot.Circle
}

// NewCustomMap initializes a new CustomMap with essential attributes
func NewCustomMap(center *geoplot.LatLng, zoom int, area *geoplot.Area) *CustomMap {
	return &CustomMap{
		Center:  center,
		Zoom:    zoom,
		Area:    area,
		Markers: []*geoplot.Marker{},
		Circles: []*geoplot.Circle{},
	}
}

// AddMarker adds a new marker to the CustomMap
func (cm *CustomMap) AddMarker(marker *geoplot.Marker) {
	cm.Markers = append(cm.Markers, marker)
}

// AddCirlce adds a new circle to the CustomMap
func (cm *CustomMap) AddCircle(circle *geoplot.Circle) {
	cm.Circles = append(cm.Circles, circle)
}

// ClearMarkers clears all markers from the CustomMap
func (cm *CustomMap) ClearMarkers() {
	cm.Markers = []*geoplot.Marker{}
	cm.Circles = []*geoplot.Circle{}
}

// ToGeoPlotMap converts the CustomMap and its markers to a geoplot.Map for rendering
func (cm *CustomMap) ToGeoPlotMap() *geoplot.Map {
	gMap := &geoplot.Map{
		Center: cm.Center,
		Zoom:   cm.Zoom,
		Area:   cm.Area,
	}
	for _, marker := range cm.Markers {
		gMap.AddMarker(marker)
	}

	for _, circle := range cm.Circles {
		gMap.AddCircle(circle)
	}
	return gMap
}

func InitHub(hubCollection map[string]*hub.Hub, hubMap *CustomMap) {
	for _, hubStation := range hubCollection {
		// Add hub to map
		icon := geoplot.ColorIcon(255, 255, 0)
		hubGeoplotCoords := &geoplot.LatLng{
			Latitude:  hubStation.Coordinates.Latitude,
			Longitude: hubStation.Coordinates.Longitude,
		}
		hubState := hubStation.ReportHubState()
		chargingSpotsAvailable := hubState[0]
		nodesWaiting := hubState[1]
		nodesAllocated := hubState[2]

		nodesCharging := hubStation.HubChargingSpotsTotal - chargingSpotsAvailable

		hubAddr := hubStation.HubHost + ":" + hubStation.HubPort

		log.Printf("Hub: %s\nCharging Spots: %d\nCharging: %d\nWaiting: %d\nAllocated: %d", hubAddr, chargingSpotsAvailable, nodesCharging, nodesWaiting, nodesAllocated)

		hubMap.AddMarker(&geoplot.Marker{
			LatLng:  hubGeoplotCoords,
			Tooltip: fmt.Sprintf("%s\nCharging Spots: %d\nCharging: %d\nWaiting: %d\nAllocated: %d", hubAddr, chargingSpotsAvailable, nodesCharging, nodesWaiting, nodesAllocated),
			Icon:    icon,
		})

		// hubMap.AddCircle(&geoplot.Circle{
		// 	LatLng:      hubGeoplotCoords,
		// 	RadiusMeter: 1000,
		// 	//Tooltip:     "Inner Circle",
		// })

		// hubMap.AddCircle(&geoplot.Circle{
		// 	LatLng:      hubGeoplotCoords,
		// 	RadiusMeter: 2000,
		// 	//Tooltip:     "Middle Circle",
		// })

		// hubMap.AddCircle(&geoplot.Circle{
		// 	LatLng:      hubGeoplotCoords,
		// 	RadiusMeter: 3000,
		// 	//Tooltip:     "Outer Circle",
		// })
	}
	log.Println("-------------------------\n\n")
}

// hubCollection is a map of hub objects
func UpdateMap(hubCollection map[string]*hub.Hub, hubMap *CustomMap, planner *planner.Planner) {

	// log.Println("Updating Map")

	hubMap.ClearMarkers()
	markers := hubMap.Markers
	InitHub(hubCollection, hubMap)

	planner.Mutex.Lock()
	defer planner.Mutex.Unlock()
	// Add nodes to map
	for nodeId, node := range planner.NodeMap {
		hubAddr := node.HubAddr

		point := &geoplot.LatLng{
			Latitude:  node.Coordinates.Latitude,
			Longitude: node.Coordinates.Longitude,
		}

		nodeIcon := geoplot.ColorIcon(0, 255, 0)
		if node.BatteryLevel <= 50 && node.BatteryLevel > 30 {
			nodeIcon = geoplot.ColorIcon(255, 128, 0)
		} else if node.BatteryLevel <= 30 {
			nodeIcon = geoplot.ColorIcon(255, 0, 0)
		}

		markers = append(markers, (&geoplot.Marker{
			LatLng: point,
			//Popup: fmt.Sprintf("NodePop: %s", nodeId),
			Tooltip: fmt.Sprintf("NodeId: %s\nBattery: %v\nSeeking Charge: %v\nHub: %v", nodeId, node.BatteryLevel, node.SeekingCharge, hubAddr),
			Icon:    nodeIcon,
		}))

	}

	// for _, hubStation := range hubCollection {
	// 	hubStation.Mutex.Lock()
	// 	defer hubStation.Mutex.Unlock()

	// 	for nodeId, node := range hubStation.NodeMap {

	// 		hubAddr := hubStation.HubHost + ":" + hubStation.HubPort

	// 		point := &geoplot.LatLng{
	// 			Latitude:  node.Coordinates.Latitude,
	// 			Longitude: node.Coordinates.Longitude,
	// 		}

	// 		nodeIcon := geoplot.ColorIcon(0, 255, 0)
	// 		if node.BatteryLevel <= 50 && node.BatteryLevel > 30 {
	// 			nodeIcon = geoplot.ColorIcon(255, 128, 0)
	// 		} else if node.BatteryLevel <= 30 {
	// 			nodeIcon = geoplot.ColorIcon(255, 0, 0)
	// 		}

	// 		markers = append(markers, (&geoplot.Marker{
	// 			LatLng: point,
	// 			//Popup: fmt.Sprintf("NodePop: %s", nodeId),
	// 			Tooltip: fmt.Sprintf("NodeId: %s\nBattery: %v\nSeeking Charge: %v\nHub: %v", nodeId, node.BatteryLevel, node.SeekingCharge, hubAddr),
	// 			Icon:    nodeIcon,
	// 		}))
	// 	}
	// }

	for _, marker := range markers {
		hubMap.AddMarker(marker)
	}
}
