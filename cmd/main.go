package main

import (
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

	// create the basic devices and remote nodes
	nodes, err := createNodes(activeNodes)

	if len(nodes) < 2 {
		log.Fatalln("Exiting 'cause not enough nodes could be initialized on given network")
	}

	n := netsimK.CreateNetwork(nodes, time.Second*5)

	n.Start()

	n.GenerateTraffic()
	go n.Wait()

	n.Stop()
}
