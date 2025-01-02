package main

import (
	"net"
	"time"
)

// type Send func(dest *Device, payload []byte) error
// type Receive func() error
type TransferCallback func(p NetPacket, src, dest Node, now time.Time)

type Node interface {
	Send(dest Node, payload []byte) error
	Receive() error

	Connect(dest Node)

	// GetTransferCallback() TransferCallback
	// SetTransferCallback(callback TransferCallback)
}

// the sore base of a node
type BaseOfNode struct {
	next     []Node
	callback TransferCallback
}

type BasicDevice struct {
	Name    string
	IP      net.IP
	Mac     net.HardwareAddr
	Active  bool
	packets chan NetPacket
}

func NewBasicDevice(name, addr string) *BasicDevice {
	return &BasicDevice{
		Name:    name,
		IP:      net.IP(addr),
		Active:  true,
		packets: make(chan NetPacket),
	}
}

func (d *BasicDevice) Connect(dest Node) {
	if device, ok := dest.(*BasicDevice); ok {
		go func() {
			for p := range device.packets {
				d.packets <- p
			}
		}()
	}
}

func (d *BasicDevice) Send(dest Node, payload []byte) error {
	return nil
}

func (d *BasicDevice) Receive() error {
	return nil
}

// func (d *BasicDevice) SetTransferCallback(callback TransferCallback) {
// 	d.callback = callback
// }
// func (d *BasicDevice) GetTransferCallback(, payload []byte) error

