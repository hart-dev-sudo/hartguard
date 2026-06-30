package sniffer

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/hart-dev-sudo/hartguard/port-scan-detector/internal/detector"
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
		if srcIP, dstPort, flags, ok := parsePacket(packet); ok {
			s.detector.Process(srcIP, dstPort, flags)
		}
	}
	return nil
}

// parsePacket extracts src IP, dst port, and TCP flags from a packet.
// Returns ok=false if the packet is not an IPv4 TCP packet.
func parsePacket(packet gopacket.Packet) (srcIP string, dstPort uint16, flags uint16, ok bool) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if ipLayer == nil || tcpLayer == nil {
		return "", 0, 0, false
	}
	ip := ipLayer.(*layers.IPv4)
	tcp := tcpLayer.(*layers.TCP)
	return ip.SrcIP.String(), uint16(tcp.DstPort), tcpFlags(tcp), true
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
