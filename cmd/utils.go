package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	netlibk "github.com/KennyZ69/netlibK"
)

const (
	MAX_NODES = 12
	BASE_PORT = 4333
)

func usage() {
	print(`
	Usage: 
		nsK <CIDR> <ifi> <port> <max_nodes>

	Please also provide root privileges.
`)
}

func getIfiFromCIDR(cidr string) (*string, error) {
	cmd := exec.Command("ip", "route", "show", cidr)
	// cidr dev ifi ......
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	lines := strings.Split(string(output), "\n") // get the returned line
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return nil, fmt.Errorf("no matching network for CIDR: %s", cidr)
	}

	fields := strings.Fields(lines[0])
	if len(fields) < 3 {
		return nil, fmt.Errorf("unexpected output format: %s", lines[0])
	}

	iface := fields[2] // ifi should be the third result in the line
	// fmt.Println(iface)
	// ifi, err := net.InterfaceByName(iface)
	return &iface, nil
}

func discoverHosts(ips []net.IP, activeHosts chan<- net.IP) error {
	var notActiveCounter, failedCounter int
	var wg sync.WaitGroup
	payload := []byte("Active? ")
	timeout := time.Second * 2

	for _, ip := range ips {

		wg.Add(1)
		go func(targetIp net.IP) {
			defer func() {
				wg.Done()
			}()
			_, active, err := netlibk.HigherLvlPing(targetIp, payload, timeout)
			if err != nil {
				failedCounter++
			}
			if active {
				// log.Printf("Host %s is active\nAdding to the list of active hosts...\n", targetIp.String())
				// log.Printf("Host %s is active with latency of %v\n", targetIp.String(), latency)
				activeHosts <- targetIp
			} else {
				// log.Printf("%s is not active host\n", targetIp.String())
				notActiveCounter++
			}
		}(ip)

	}

	go func() {
		wg.Wait() // Ensure all goroutines finish
		// Close channel after all pings are done
		close(activeHosts)
	}()

	var err error = nil
	if failedCounter == len(ips) {
		err = fmt.Errorf("All tries for pings failed, you may need to run this with sudo or there is another problem... check your reports logs")
	}

	fmt.Println("If some host you expected to be active seems to not be, you may need to run this tool with sudo (admin privileges)")

	return err
}

func getInputIPs(ifaceFlag, ipStart *string) ([]net.IP, *net.Interface) {
	var addrs []string
	var ipArr []net.IP
	ifi, err := net.InterfaceByName(*ifaceFlag)
	if err != nil {
		log.Fatalf("Error getting the net interface: %v\n", err)
	}

	addrs = append(addrs, *ipStart, "")

	addr_start, addr_end, isCidr, err := netlibk.ParseIPInputs(addrs)
	if err != nil {
		log.Fatalf("Error parsing the input values: %v\n", err)
	}

	if isCidr {
		ipArr = netlibk.GenerateIPsFromCIDR(addr_start)
	} else {
		ipArr = netlibk.GenerateIPs(addr_start, addr_end)
	}

	return ipArr, ifi
}
