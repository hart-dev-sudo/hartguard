package sniffer

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func makePacket(srcIP string, dstPort uint16, syn, fin, psh, urg bool) gopacket.Packet {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}

	ip := &layers.IPv4{
		SrcIP:    net.ParseIP(srcIP).To4(),
		DstIP:    net.ParseIP("10.0.0.1").To4(),
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		TTL:      64,
	}
	tcp := &layers.TCP{
		SrcPort: 12345,
		DstPort: layers.TCPPort(dstPort),
		SYN:     syn,
		FIN:     fin,
		PSH:     psh,
		URG:     urg,
	}
	tcp.SetNetworkLayerForChecksum(ip)
	gopacket.SerializeLayers(buf, opts, ip, tcp)
	return gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
}

func TestParsePacketSYN(t *testing.T) {
	pkt := makePacket("192.168.1.50", 80, true, false, false, false)
	srcIP, dstPort, flags, ok := parsePacket(pkt)
	if !ok {
		t.Fatal("expected ok=true for valid TCP/IP packet")
	}
	if srcIP != "192.168.1.50" {
		t.Errorf("expected srcIP 192.168.1.50, got %s", srcIP)
	}
	if dstPort != 80 {
		t.Errorf("expected dstPort 80, got %d", dstPort)
	}
	if flags != 0x002 {
		t.Errorf("expected SYN flags 0x002, got 0x%03x", flags)
	}
}

func TestParsePacketXMAS(t *testing.T) {
	pkt := makePacket("10.0.0.5", 443, false, true, true, true)
	_, _, flags, ok := parsePacket(pkt)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if flags != 0x029 {
		t.Errorf("expected XMAS flags 0x029, got 0x%03x", flags)
	}
}

func TestParsePacketNULL(t *testing.T) {
	pkt := makePacket("10.0.0.6", 22, false, false, false, false)
	_, _, flags, ok := parsePacket(pkt)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if flags != 0x000 {
		t.Errorf("expected NULL flags 0x000, got 0x%03x", flags)
	}
}

func TestTCPFlagsSYN(t *testing.T) {
	tcp := &layers.TCP{SYN: true}
	if f := tcpFlags(tcp); f != 0x002 {
		t.Errorf("expected 0x002, got 0x%03x", f)
	}
}

func TestTCPFlagsFIN(t *testing.T) {
	tcp := &layers.TCP{FIN: true}
	if f := tcpFlags(tcp); f != 0x001 {
		t.Errorf("expected 0x001, got 0x%03x", f)
	}
}
