package widgets

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
)

const NetworkWidgetType = "network"

type NetworkWidget struct {
	Interface   string `json:"interface"`
	Interval    int    `json:"interval"`
	Received    int    `json:"received"`
	Transmitted int    `json:"transmitted"`
}

func (widget *NetworkWidget) Name() string {
	if widget.Interface != "" {
		return widget.Interface
	}
	return "Network"
}

func (widget *NetworkWidget) Type() string {
	return "network"
}

func (widget *NetworkWidget) HasData() bool {
	return true
}

func (widget *NetworkWidget) Update() error {
	return errors.New("Not implemented")
}

type NetworkDevice struct {
	Interface        string
	ReceivedBytes    uint64
	TransmittedBytes uint64
}

func GatherNetDevInfo() ([]NetworkDevice, error) {
	fh, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	result := make([]NetworkDevice, 0)
	reader := bufio.NewReader(fh)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		dev := NewNetworkDevice(line)
		if dev != nil {
			result = append(result, *dev)
		}
	}

	return result, nil
}

var NetworkDeviceSplitRegexp = regexp.MustCompile("\\s+")

func NewNetworkDevice(line string) *NetworkDevice {
	parts := NetworkDeviceSplitRegexp.Split(line, -1)
	if len(parts) == 0 {
		return nil
	}
	if parts[0] == "" {
		parts = parts[1:]
	}
	if len(parts) < 10 {
		return nil
	}

	iface := parts[0]
	if len(iface) < 0 || iface[len(iface)-1] != ':' {
		return nil
	}

	iface = iface[:len(iface)-1]
	recvd, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil
	}
	trans, err := strconv.ParseUint(parts[9], 10, 64)
	if err != nil {
		return nil
	}

	return &NetworkDevice{
		Interface:        iface,
		ReceivedBytes:    recvd,
		TransmittedBytes: trans,
	}
}
