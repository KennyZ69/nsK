package netsimK

import "sync"

type PacketRouter struct {
	nodes []Node
	// nodes map[string]*Node
	mu sync.Mutex
}

func NewPacketRouter() *PacketRouter {
	// return &PacketRouter{nodes: map[string]*Node{}, mu: sync.Mutex{}}
	return &PacketRouter{nodes: []Node{}, mu: sync.Mutex{}}
}

func (r *PacketRouter) AddNode(node Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes = append(r.nodes, node)
}

func (r *PacketRouter) SendPacket(p NetPacket) {

}
