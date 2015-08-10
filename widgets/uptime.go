package widgets

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"
)

const UptimeWidgetType = "uptime"

type UptimeWidget struct {
	LastUpdated time.Time
	Days        int
	Hours       int
	Minutes     int
	Seconds     int
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

func (widget *UptimeWidget) UnmarshalJSON(data []byte) error {
	var helper struct {
		Days    int `json:"days"`
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
		Seconds int `json:"seconds"`
	}

	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}

	widget.Days = helper.Days
	widget.Hours = helper.Hours
	widget.Minutes = helper.Minutes
	widget.Seconds = helper.Seconds
	widget.LastUpdated = time.Now()

	return nil
}

func (widget *UptimeWidget) MarshalJSON() ([]byte, error) {
	var helper struct {
		Days    int `json:"days"`
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
		Seconds int `json:"seconds"`
	}

	total := CombineToSeconds(
		widget.Days,
		widget.Hours,
		widget.Minutes,
		widget.Seconds,
	)

	if !widget.LastUpdated.IsZero() {
		total += int(time.Since(widget.LastUpdated).Seconds())
	}

	d, h, m, s := SplitSeconds(total)

	helper.Days = d
	helper.Hours = h
	helper.Minutes = m
	helper.Seconds = s

	return json.Marshal(helper)
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

	days, hours, minutes, seconds = SplitSeconds(seconds)

	return
}

func SplitSeconds(in int) (days, hours, minutes, seconds int) {
	seconds = in

	days = seconds / (60 * 60 * 24)
	seconds -= days * (60 * 60 * 24)

	hours = seconds / (60 * 60)
	seconds -= hours * (60 * 60)

	minutes = seconds / 60
	seconds -= minutes * 60

	return
}

func CombineToSeconds(days, hours, minutes, seconds int) int {
	hours += (days * 24)
	minutes += (hours * 60)
	seconds += (minutes * 60)

	return seconds
}
