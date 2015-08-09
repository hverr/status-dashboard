package widgets

import (
	"bufio"
	"errors"
	"os"
	"strconv"
)

const UptimeWidgetType = "uptime"

type UptimeWidget struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

func (widget *UptimeWidget) Name() string {
	return "Uptime"
}

func (widget *UptimeWidget) Type() string {
	return UptimeWidgetType
}

func (widget *UptimeWidget) HasData() bool {
	return true
}

func (widget *UptimeWidget) Update() error {
	d, h, m, s, err := GatherUptime()
	if err != nil {
		return err
	}

	widget.Days, widget.Hours = d, h
	widget.Minutes, widget.Seconds = m, s

	return nil
}

func GatherUptime() (days, hours, minutes, seconds int, err error) {
	fh, err := os.Open("/proc/uptime")
	if err != nil {
		return
	}
	defer fh.Close()

	reader := bufio.NewReader(fh)
	line, err := reader.ReadString('.')
	if err != nil {
		return
	} else if len(line) <= 1 {
		err = errors.New("Unexpected /proc/uptime contents.")
		return
	}

	seconds, err = strconv.Atoi(line[:len(line)-1])
	if err != nil {
		err = errors.New("Could not parse /proc/uptime: " + err.Error())
		return
	}

	days = seconds / (60 * 60 * 24)
	seconds -= days * (60 * 60 * 24)

	hours = seconds / (60 * 60)
	seconds -= hours * (60 * 60)

	minutes = seconds / 60
	seconds -= minutes * 60

	return
}
