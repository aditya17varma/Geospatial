# GeoSpatial Tracking and Resource Allocation

## Introduction
A project to find efficient resource allocation policies for moving nodes.

We study the policies to allocate charging hubs to drones. We assume there are 2 Charging Hubs. Each hub will have charging spots, the number of which can be set in the config.toml file.
Some policies studied are: "minDistance", "minBattery", "EM".
Efficient Allocation is determined by the Total and Average downtime a Node spends waiting to charge at one of the Hubs.
The lower the downtime, the less time Nodes are spending waiting.

1. **minDistance**: The Hub closest to the drone is chosen.
2. **minBattery**: The Hub closest to the drone is chosen, but at the drone itself, the drones with the lowest battery level are given preference.
3. **EM**: [**Expectation Maximization**](https://en.wikipedia.org/wiki/Expectation–maximization_algorithm). An ML algorithm. A cost function that takes in the distance, expected wait time, and battery level is used to select a Hub. This is a dynamic function that uses the everchanging states of the Hubs and Nodes to find the best allocation in momemt of time.


## Setup and Run
The config.toml file lets you choose the Allocation Policy, the number of Charging Spots per Hub, the Number of Drones/Nodes, as well as the simulation Time.

To run:
```go run main.go```


## Project Architecture
There are 3 main elements to the Project:
1. Planner
2. Hub
3. Node

All communication uses Protobufs. The messages are defined in the proto/heartbeat.proto file.

### Planner
The planner is responsible for allocating the resource Hub to the Nodes. It performs the relevant calculation and assigns the Hub.
All nodes register with the Planner at the start.
When the nodes are seeeking charging, they communicate with the Planner and the Planner resonds with the assigned Hub information.
The Planner receives state updates in the form of heartbeats from the Nodes.

### Hub
Hubs are where the Nodes recive the resources / charging.
Once the Planner, assigns the Hub, the Nodes move towards the Hub.
Once a Node reaches the Hub and charging spots are available, the Node takes up a charging spot and begins charging.
If no spots are open, the Node joins the waiting queue. Nodes are taken off the queue according to the allocation policy.
The Hub receives state updates in the form of heartbeats from the Nodes once they are assigned to the Hub.

### Node
Nodes are either moving towards a target (randomly allocated),a Hub (allocated by Planner), or charging (at Hub).
Nodes (or drones) are moving objects. They need to keep track of their geospatial coordinates, and state to communicate effectively.
Nodes send state updates in the form of heartbeat messages to the Planner and Hubs.

## Findings
Simulations were run for each Allocation Policy for 100 seconds and 30 nodes.
Here are the findings for the "minDistance", "minBattery", and "EM" (Expectation Maximization) policies:

**minDistance**

![image](https://github.com/aditya17varma/Geospatial/assets/69626061/db32b2f9-f87a-46af-a3df-8de8468c87cf)


**minBattery**

![image](https://github.com/aditya17varma/Geospatial/assets/69626061/70334bf3-2962-44a1-a275-82eead09a427)


**Expecation Maximization**

![image](https://github.com/aditya17varma/Geospatial/assets/69626061/f01e61ba-9020-4c37-a288-5e35dd8917d9)


As we can see EM is the best allocation policy, with a significant reduction in Total and Average Downtime.


## Challenges
Some of the challenges I faced when building this project:
<ul>
  <li>Planning the Proper Architecutre. With so many moving pieces, it was important to find an architeture that didn't have inefficient communication.</li>
  <li>Debugging Multi-Threaded Processes. Each channel of communication with each node was it's own go routine. Multiple go routines were also used to keep track of the simulation state.</li>
  <li>Visualization. Go doesn't have a visualization library as powerful and expansive as Python. In the future, I might pipe the data from the simulation to a Python/React server that displays realtime movement of the nodes and other stats.</li>
</ul>


