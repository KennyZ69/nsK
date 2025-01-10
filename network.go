package netsimK

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Network struct {
	nodes    []Node
	wg       *sync.WaitGroup
	lifetime time.Duration
}

func (n *Network) Start() {
	log.Println("Starting up the network ... ")

	for _, node := range n.nodes {
		node.Start()
	}
}

func (n *Network) Stop() {
	log.Println("Stopping the network ... ")
	for _, node := range n.nodes {
		node.Close()
	}
}

func CreateNetwork(nodes []Node, lifetime time.Duration) *Network {
	log.Println("Creating a network from", len(nodes), "number of nodes")
	return &Network{
		nodes:    nodes,
		wg:       &sync.WaitGroup{},
		lifetime: lifetime,
	}
}

// Wait until network simulation finishes
func (n *Network) Wait() {
	n.wg.Wait()
}

func (n *Network) AddNode(node Node) {
	n.nodes = append(n.nodes, node)
}

func (n *Network) GenerateTraffic() {
	log.Println("Generating traffic ...")
	rep := make(chan string)
	n.wg.Add(1)
	for _, src := range n.nodes {
		for _, dest := range n.nodes {
			if src != dest {
				// go src.Connect(dest)
				go func(src, dest Node) {
					defer n.wg.Done()
					for {
						p := &SimPacket{
							// Source:  src,
							// Dest:    dest,
							Payload: []byte("Testing payload!"),
						}
						if err := src.Send(dest, p.Payload); err != nil {
							log.Printf("[%v] failed to send packet to [%v] -> %v\n", src, dest, err)
							rep <- fmt.Sprintf("[%v] failed to send packet to [%v] -> %v\n", src, dest, err)
						} else {
							log.Printf("Successfully sent a packet from [%v] to [%v]\n", src, dest)
							// time.Sleep(1 * time.Second)
							rep <- fmt.Sprintf("Successfully sent a packet from [%v] to [%v]\n", src, dest)
						}
					}
				}(src, dest)
				go func() {
					select {
					case msg := <-rep:
						log.Print(msg)
					}
				}()
			}
		}
	}
}
