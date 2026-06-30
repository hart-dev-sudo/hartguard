package sniffer

import (
	"fmt"
	"log"

	"github.com/hart-dev-sudo/hartguard/port-scan-detector/internal/detector"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Sniffer struct {
	iface    string
	detector *detector.Detector
}

func New(iface string, d *detector.Detector) *Sniffer {
	return &Sniffer{iface: iface, detector: d}
}

func (s *Sniffer) Start() error {
	handle, err := pcap.OpenLive(s.iface, 1600, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("opening interface %s: %w", s.iface, err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("tcp"); err != nil {
		return fmt.Errorf("setting BPF filter: %w", err)
	}

	log.Printf("Listening on %s", s.iface)

	src := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range src.Packets() {
		s.processPacket(packet)
	}
	return nil
}

func (s *Sniffer) processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if ipLayer == nil || tcpLayer == nil {
		return
	}

	ip := ipLayer.(*layers.IPv4)
	tcp := tcpLayer.(*layers.TCP)

	s.detector.Process(ip.SrcIP.String(), uint16(tcp.DstPort), tcpFlags(tcp))
}

func tcpFlags(tcp *layers.TCP) uint16 {
	var f uint16
	if tcp.SYN {
		f |= 0x002
	}
	if tcp.FIN {
		f |= 0x001
	}
	if tcp.PSH {
		f |= 0x008
	}
	if tcp.URG {
		f |= 0x020
	}
	return f
}
