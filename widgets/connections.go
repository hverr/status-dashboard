package widgets

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
)

const ConnectionsWidgetType = "connections"

type ConnectionsWidget struct {
	TCP4 int `json:"tcp4"`
	TCP6 int `json:"tcp6"`
}

func (widget *ConnectionsWidget) Name() string {
	return "Connections"
}

func (widget *ConnectionsWidget) Type() string {
	return ConnectionsWidgetType
}

func (widget *ConnectionsWidget) HasData() bool {
	return true
}

func (widget *ConnectionsWidget) Configure(json.RawMessage) error {
	return nil
}

func (widget *ConnectionsWidget) Configuration() interface{} {
	return nil
}

func (widget *ConnectionsWidget) Start() error {
	return nil
}

func (widget *ConnectionsWidget) Update() error {
	tcp4, tcp6, err := GatherTCPConnections()
	if err != nil {
		return err
	}

	widget.TCP4 = tcp4
	widget.TCP6 = tcp6

	return nil
}

type TCPState int

const (
	TCPEstablished TCPState = 1
	TCPSynSent     TCPState = 2
	TCPSynRecv     TCPState = 3
	TCPFinWait1    TCPState = 4
	TCPFinWait2    TCPState = 5
	TCPTimeWait    TCPState = 6
	TCPClose       TCPState = 7
	TCPCloseWait   TCPState = 8
	TCPLastAck     TCPState = 9
	TCPListen      TCPState = 10
	TCPClosing     TCPState = 11
)

type TCPInfo struct {
	Entry      uint     // sl
	LocalIP    net.IP   // local_address
	LocalPort  uint16   // local_address
	RemoteIP   net.IP   // remote_address
	RemotePort uint16   // remote_address
	State      TCPState // st
}

var TCPInfoRegexp = regexp.MustCompile("^\\s*(\\d+): ([0-9a-fA-F]+):([0-9a-fA-F]+) ([0-9a-fA-F]+):([0-9a-fA-F]+) ([0-9a-fA-F]+)")

func NewTCPInfo(line string) (*TCPInfo, error) {
	m := TCPInfoRegexp.FindStringSubmatch(line)
	if m == nil {
		return nil, errors.New("Line did not match TCPInfoRegexp")
	}

	i := TCPInfo{}
	if ui, err := strconv.ParseUint(m[1], 10, 32); err == nil {
		i.Entry = uint(ui)
	}

	if ui, err := strconv.ParseUint(m[3], 16, 16); err == nil {
		i.LocalPort = uint16(ui)
	}

	if ui, err := strconv.ParseUint(m[5], 16, 16); err == nil {
		i.RemotePort = uint16(ui)
	}

	if ui, err := strconv.ParseInt(m[6], 16, 32); err == nil {
		i.State = TCPState(int(ui))
	}

	localBytes, err := hex.DecodeString(m[2])
	if err != nil {
		return nil, err
	}
	i.LocalIP = bytesToIP(localBytes)

	remoteBytes, err := hex.DecodeString(m[4])
	if err != nil {
		return nil, err
	}
	i.RemoteIP = bytesToIP(remoteBytes)

	return &i, nil
}

func GatherTCPConnections() (tcp4, tcp6 int, err error) {

	infos4, err := GatherTCP4Infos()
	if err != nil {
		return
	}

	infos6, err := GatherTCP6Infos()
	if err != nil {
		return
	}

	local, err := GatherLocalIPAddresses()
	if err != nil {
		return
	}

	isLocal := func(ip net.IP) bool {
		for _, l := range local {
			if l.Equal(ip) {
				return true
			}
		}
		return false
	}

	for _, j := range infos4 {
		if j.State == TCPEstablished && (!isLocal(j.LocalIP) || !isLocal(j.RemoteIP)) {
			tcp4 += 1
		}
	}

	for _, j := range infos6 {
		if j.State == TCPEstablished && (!isLocal(j.LocalIP) || !isLocal(j.RemoteIP)) {
			tcp6 += 1
		}
	}

	return
}

func GatherTCPInfos(filename string) ([]TCPInfo, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	var result []TCPInfo

	reader := bufio.NewReader(fh)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		i, _ := NewTCPInfo(line)
		if i != nil {
			result = append(result, *i)
		}
	}
	return result, nil
}

func GatherTCP4Infos() ([]TCPInfo, error) {
	return GatherTCPInfos("/proc/net/tcp")
}

func GatherTCP6Infos() ([]TCPInfo, error) {
	return GatherTCPInfos("/proc/net/tcp6")
}

func GatherLocalIPAddresses() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ips = append(ips, v.IP)
			case *net.IPAddr:
				ips = append(ips, v.IP)
			}
		}
	}

	return ips, nil
}

func bytesToIP(bytes []byte) net.IP {
	switch len(bytes) {
	case 4:
		reverseIPBytes(bytes[0:4])
	case 16:
		reverseIPBytes(bytes[0:4])
		reverseIPBytes(bytes[4:8])
		reverseIPBytes(bytes[8:12])
		reverseIPBytes(bytes[12:16])
	default:
		return nil
	}

	return net.IP(bytes)
}

func reverseIPBytes(bytes []byte) {
	j := len(bytes)
	m := j / 2
	j -= 1
	i := 0
	for i < m {
		b := bytes[i]
		bytes[i] = bytes[j]
		bytes[j] = b

		i += 1
		j -= 1
	}
}
