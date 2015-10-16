package widgets

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/pmylund/go-cache"
)

const NetworkWidgetType = "network"

type NetworkWidget struct {
	Interface   string  `json:"interface"`
	Interval    float64 `json:"interval"`
	Received    int     `json:"received"`
	Transmitted int     `json:"transmitted"`

	configuration struct {
		Interface string `json:"interface"`
	} `json:"configuration"`

	hasData bool `json:"-"`
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

func (widget *NetworkWidget) Identifier() string {
	return widget.Type() + "_" + widget.configuration.Interface
}

func (widget *NetworkWidget) HasData() bool {
	return widget.hasData
}

func (widget *NetworkWidget) Configure(c json.RawMessage) error {
	if err := json.Unmarshal(c, &widget.configuration); err != nil {
		return err
	}

	widget.Interface = widget.configuration.Interface

	return nil
}

func (widget *NetworkWidget) Configuration() interface{} {
	return widget.configuration
}

var updateProcess = networkInformationProvider{}

func (widget *NetworkWidget) Start() error {
	updateProcess.Interval = 1 * time.Second
	updateProcess.Start()
	return nil
}

func (widget *NetworkWidget) Update() error {
	if !updateProcess.Started() {
		return errors.New("Network widget was not started")
	}

	w, ok := updateProcess.WidgetForInterface(widget.Interface)
	if !ok {
		widget.hasData = false
		return nil
	}

	widget.Interval = w.Interval
	widget.Received = w.Received
	widget.Transmitted = w.Transmitted

	widget.hasData = true

	return nil
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

// networkInformationProvider is a structure to hold information
// about all interfaces of the system.
type networkInformationProvider struct {
	Interval time.Duration

	dispatcher sync.Once
	widgets    *cache.Cache
	stopper    chan bool

	previousTime    time.Time
	previousDevices []NetworkDevice
}

func (p *networkInformationProvider) Start() {
	go p.dispatcher.Do(func() {
		p.stopper = make(chan bool)
		p.widgets = cache.New(cache.NoExpiration, cache.NoExpiration)

		p.run()
	})
}

func (p *networkInformationProvider) Started() bool {
	return p.stopper != nil
}

func (p *networkInformationProvider) Stop() {
	if p.stopper != nil {
		p.stopper <- true
	}
}

func (p *networkInformationProvider) WidgetForInterface(iface string) (NetworkWidget, bool) {
	o, found := p.widgets.Get(iface)
	if !found {
		return NetworkWidget{}, false
	}

	return o.(NetworkWidget), true
}

func (p *networkInformationProvider) run() {
	stopped := false
	for !stopped {
		select {
		case <-time.After(p.Interval):
			p.updateAll()
		case <-p.stopper:
			stopped = true
		}
	}
}

func (provider *networkInformationProvider) updateAll() {
	now := time.Now()
	devices, err := GatherNetDevInfo()
	if err != nil {
		log.Println("Could not gather network information:", err)
		provider.widgets.Flush()
		return
	}

	// Clear widgets.
	provider.widgets.Flush()

	// Calculate average
	if provider.previousDevices != nil {
		duration := now.Sub(provider.previousTime)

		for _, p := range provider.previousDevices {
			for _, d := range devices {
				if p.Interface == d.Interface {
					w := NetworkWidget{
						Interface:   d.Interface,
						Interval:    duration.Seconds(),
						Received:    int(d.ReceivedBytes - p.ReceivedBytes),
						Transmitted: int(d.TransmittedBytes - p.TransmittedBytes),
					}
					provider.widgets.Set(d.Interface, w, cache.DefaultExpiration)
				}
			}
		}
	}

	// Set previous
	provider.previousDevices = devices
	provider.previousTime = now
}
