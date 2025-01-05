package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/KennyZ69/netsimGo"
)

func main() {
	// a user should give me an IP (CIDR) range and a port that I could use for the generated devices
	// from the ip range I would discover active hosts and use them as my nodes
	// and third arg could be desired max number of nodes for this sim network

	args := os.Args[1:]

	iface, err := getIfiFromCIDR(args[0])
	if err != nil || iface == nil {
		log.Fatalf("Couldn't fetch network interface, please provide it yourself: %s: %v\n", *iface, err)
	}

	ips, _ := getInputIPs(iface, &args[0])
	activeNodes := make(chan net.IP)
	if err = discoverHosts(ips, activeNodes); err != nil {
		log.Fatal(err)
	}

	// if len(activeNodes) == 0 {
	// 	log.Fatalln("Error: There were no active nodes found on given network")
	// }

	var nodes []netsimK.Node
	var i, port int = 0, BASE_PORT
	for n := range activeNodes {
		log.Println("Host", n.String(), "is active")
		d, err := netsimK.NewBasicDevice(fmt.Sprintf("D%d", i), n.String(), port)
		if err != nil || d == nil {
			continue // I will skip this device
		}
		nodes = append(nodes, d)
		i++
		port++
	}

	if len(nodes) < 2 {
		log.Fatalln("Exiting 'cause not enough nodes could be initialized on given network")
	}

	n := netsimK.CreateNetwork(nodes, time.Second*5)

	n.Start()
	defer n.Stop()

	go n.GenerateTraffic()

	// d1 := netsimK.NewBasicDevice("D1", "192.168.0.1", 3900)
	// d2 := netsimK.NewBasicDevice("D2", "192.168.0.2", 3901)
	// d3 := netsimK.NewBasicDevice("D3", "192.168.0.3", 3902)
	// nodes := []netsimK.Node{d1, d2, d3}
}
