package main

import (
	"sync"
	"time"
)

type Network struct {
	nodes    []*Node
	wg       *sync.WaitGroup
	lifetime time.Duration
	// buf
}

func (n *Network) Run() {

}

func CreateNetwork(nodes []Node, lifetime time.Duration) *Network

// Wait until network simulation finishes
func (n *Network) Wait() {
	n.wg.Wait()
}
