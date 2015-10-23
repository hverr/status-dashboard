package widgets

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

const MeminfoWidgetType = "meminfo"

type MeminfoWidget struct {
	Total int `json:"total"`
	Used  int `json:"used"`
}

func (widget *MeminfoWidget) Name() string {
	return "Widget"
}

func (widget *MeminfoWidget) Type() string {
	return MeminfoWidgetType
}

func (widget *MeminfoWidget) Identifier() string {
	return widget.Type()
}

func (widget *MeminfoWidget) HasData() bool {
	return true
}

func (widget *MeminfoWidget) Configure(json.RawMessage) error {
	return nil
}

func (widget *MeminfoWidget) Configuration() interface{} {
	return nil
}

func (widget *MeminfoWidget) Start() error {
	return nil
}

func (widget *MeminfoWidget) Update() error {
	m, err := GatherMeminfo()
	if err != nil {
		return err
	}

	buffersAndCache := m.MainCachedKB + m.MainBuffersKB

	widget.Total = int(m.MainTotalKB >> 10)
	widget.Used = int(m.MainUsedKB-buffersAndCache) >> 10

	return nil
}

type Meminfo struct {
	MainTotalKB   uint64
	MainFreeKB    uint64
	MainBuffersKB uint64
	MainCachedKB  uint64

	MainUsedKB uint64
}

var MeminfoRegexp = regexp.MustCompile("^([^:]*):\\s+(\\d+)")

func GatherMeminfo() (*Meminfo, error) {
	// Code largely taken from
	//  - procps/free.c
	//  - procps/proc/sysinfo.c

	fh, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	m := &Meminfo{}

	reader := bufio.NewReader(fh)
	numFields := 4
	matched := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		matches := MeminfoRegexp.FindStringSubmatch(line)
		if matches == nil {
			err = errors.New("Line in /proc/meminfo was unexpected: " + line)
			return nil, err
		}

		field := matches[1]
		kb, err := strconv.ParseUint(matches[2], 10, 64)
		if err != nil {
			return nil, err
		}

		switch field {
		case "MemTotal":
			m.MainTotalKB = kb
			matched++
		case "MemFree":
			m.MainFreeKB = kb
			matched++
		case "Buffers":
			m.MainBuffersKB = kb
			matched++
		case "Cached":
			m.MainCachedKB = kb
			matched++
		}
	}

	if matched != numFields {
		err = errors.New(fmt.Sprintf("Could not fill meminfo due to missing fields (%d/%d filled)", matched, numFields))
		return nil, err
	}

	m.MainUsedKB = m.MainTotalKB - m.MainFreeKB

	return m, nil
}
