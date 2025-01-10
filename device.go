package netsimK

import (
	"fmt"
	"log"
	"net"
	"time"
)

// type Send func(dest *Device, payload []byte) error
// type Receive func() error
// type TransferCallback func(p NetPacket, src, dest Node, now time.Time)

type Node interface {
	Send(dest Node, payload []byte) error
	Receive(packet NetPacket) error

	// Connect(dest Node)

	Start()
	Close()

	// GetTransferCallback() TransferCallback
	// SetTransferCallback(callback TransferCallback)
}

type BasicDevice struct {
	Name     string
	IP       net.IP
	Mac      net.HardwareAddr
	Port     int // listening port
	Active   bool
	Listener net.Listener
	packets  chan NetPacket

	// callback TransferCallback
}

type RemoteHost struct {
	Name    string
	IP      net.IP
	Port    int
	Active  bool
	packets chan NetPacket
}

func (h *RemoteHost) Send(dest Node, paylod []byte) error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", h.IP.String(), h.Port))
	if err != nil {
		return fmt.Errorf("[%s] Error connecting to remote node: %v\n", h.Name, err)
	}
	defer conn.Close()

	_, err = conn.Write(paylod)
	if err != nil {
		return fmt.Errorf("[%s] Error sending a packet: %v\n", h.Name, err)
	}
	return nil
}
func (h *RemoteHost) Start() {
	log.Printf("[%s] starting up ... \n", h.Name)
	// go func() {
	// 	for h.Active {
	// 		go h.Read()
	// 	}
	// }()
}
func (h *RemoteHost) Close() {
	log.Printf("[%s] closing down ... \n", h.Name)
	h.Active = false
}

func (h *RemoteHost) Receive(packet NetPacket) error {
	log.Printf("[%s] received a net packet\n", h.Name)
	switch packet.(type) {
	case *SimPacket:
		if !packet.(*SimPacket).Ack {
			h.packets <- packet
			p := SimPacket{
				Source:  h,
				Dest:    packet.(*SimPacket).Source,
				Ack:     true,
				Payload: []byte("ACK"),
			}
			h.Send(h.router, p.Marshall()) // maybe I could start using the router
			// Because it should be a node also so that could work
		}
	}
	return nil
}
func NewRemoteNode(name string, ip net.IP, port int) *RemoteHost {
	log.Printf("Node [%s] - %s:%d -> succesfully initialized\n", name, ip.String(), port)
	return &RemoteHost{
		Name:    name,
		IP:      ip,
		Port:    port,
		Active:  true,
		packets: make(chan NetPacket, 50),
	}
}

func NewBasicDevice(name, addr string, port int) (*BasicDevice, error) {
	l, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", addr, port))
	// l, err := net.Listen("tcp4", fmt.Sprintf("%s", addr))
	if err != nil {
		// log.Printf("Failed to initialize device on %s:%d -> %v\n", addr, port, err)
		return nil, err
	}
	log.Printf("Node [%s] - %s:%d -> succesfully initialized\n", name, addr, port)
	return &BasicDevice{
		Name:     name,
		IP:       net.IP(addr),
		Active:   true,
		Port:     port,
		packets:  make(chan NetPacket, 50),
		Listener: l,
	}, nil
}

func (d *BasicDevice) Start() {
	log.Printf("[%s] starting up ...\n", d.Name)
	go func() {

		for d.Active { // 'till it is active
			// defer d.Listener.Close()
			conn, err := d.Listener.Accept()
			if err != nil {
				log.Printf("[%s] failed to accept connection -> %v\n", d.Name, err)
				// return
				continue
			}
			go d.Read(conn)
		}
		d.Listener.Close()
	}()
}

func (d *BasicDevice) Close() {
	log.Printf("[%s] closing down \n", d.Name)
	d.Active = false
	// d.Listener.Close()
}

func (d *BasicDevice) Read(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("[%s] failed to read from connection -> %v\n", d.Name, err)
			break
		}
		log.Printf("[%s] received data: %s\n", d.Name, string(buf[:n]))
	}
}

func (d *BasicDevice) Connect(dest Node) {
	log.Printf("[%s] connecting to [%s]\n", d.Name, dest.(*BasicDevice).Name)
	if device, ok := dest.(*BasicDevice); ok {
		go func() {
			for range device.packets {
				// d.packets <- p
				d.Send(device, []byte("Hello world!"))
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func (d *BasicDevice) Send(dest Node, payload []byte) error { // payload could already be a marshalled packet
	switch dest.(type) {
	case *BasicDevice:
		log.Printf("[%s] sending packet to [%s]: %s\n", d.Name, dest.(*BasicDevice).Name, string(payload))
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dest.(*BasicDevice).IP.String(), dest.(*BasicDevice).Port))
		if err != nil {
			// log.Printf("[%s] failed to connect to [%s] -> %v", d.Name, dest.(*BasicDevice).Name, err)
			return err
		}
		defer conn.Close()
		_, err = conn.Write(payload)
		if err != nil {
			return err
		}
	case *RemoteHost:
		log.Printf("[%s] sending packet to [%s]: %s\n", d.Name, dest.(*RemoteHost).Name, string(payload))
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dest.(*RemoteHost).IP.String(), dest.(*RemoteHost).Port))
		if err != nil {
			// log.Printf("[%s] failed to connect to [%s] -> %v", d.Name, dest.(*BasicDevice).Name, err)
			return err
		}
		defer conn.Close()
		_, err = conn.Write(payload)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("Shouldn't have gotten to this line\n")
}

func (d *BasicDevice) Receive(p NetPacket) error {
	// for p := range d.packets {
	// 	log.Printf("[%s] received packet from [%s]: %s", d.Name, p.(*SimPacket).Source.(*BasicDevice).Name, string(p.(*SimPacket).Payload))
	// 	if d.callback != nil {
	// 		d.callback(p, p.(*SimPacket).Source, p.(*SimPacket).Dest, time.Now())
	// 	}
	// }
	return nil
}

// func (d *BasicDevice) SetTransferCallback(callback TransferCallback) {
// 	d.callback = callback
// }
// func (d *BasicDevice) GetTransferCallback(, payload []byte) error
