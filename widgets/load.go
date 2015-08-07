package widgets

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
)

type LoadWidget struct {
	One     string `json:"one"`
	Five    string `json:"five"`
	Fifteen string `json:"fifteen"`
	Cores   int    `json:"cores"`
}

func (widget LoadWidget) Name() string {
	return "Load"
}

func (widget LoadWidget) Type() string {
	return "load"
}

func (widget LoadWidget) HasData() bool {
	return true
}

func (widget *LoadWidget) GatherInformation() error {
	if widget.Cores == 0 {
		i, err := GatherCores()
		if err != nil {
			return err
		}
		widget.Cores = i
	}

	one, five, fifteen, err := GatherLoadAverage()
	if err != nil {
		return err
	}

	widget.One, widget.Five, widget.Fifteen = one, five, fifteen

	return nil
}

var LoadCoresRegexp = regexp.MustCompile("^processor")

func GatherCores() (int, error) {
	fh, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}
	defer fh.Close()

	counter := 0
	reader := bufio.NewReader(fh)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		if LoadCoresRegexp.MatchString(line) {
			counter += 1
		}
	}

	return counter, nil
}

var LoadAverageRegexp = regexp.MustCompile("^(.*?)\\s+(.*?)\\s+(.*?)")

func GatherLoadAverage() (one, five, fifteen string, err error) {
	fh, err := os.Open("/proc/loadavg")
	if err != nil {
		return
	}

	reader := bufio.NewReader(fh)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	matches := LoadAverageRegexp.FindStringSubmatch(line)
	if matches == nil {
		err = errors.New("First line of /proc/loadavg did not match regexp.")
		return
	}

	one = matches[1]
	five = matches[2]
	fifteen = matches[3]

	return
}
