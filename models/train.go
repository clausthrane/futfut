package models

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var numberSequenceRegex *regexp.Regexp

func init() {
	numberSequenceRegex, _ = regexp.Compile("([0-9]+)")
}

// Train respresents a train event
//
// E.g. when does a certain train (by TrainNumber) arrive at station Station Uic
type Train struct {
	ID                     string `UUID`
	StationUic             string `int`
	TrainType              string
	DestinationName        string
	DestinationID          int
	Track                  string `int`
	Generated              string `date`
	TrainNumber            string `int`
	DestinationCountryCode string `int`
	ScheduledArrival       string `date`
	ScheduledDeparture     string `date`
	ArrivalDelay           string `int`
	DepartureDelay         string `int`
	Cancelled              bool
	Line                   string
	Direction              string `Nord : North, Syd: South, ..`
	MinutesToDeparture     string `int`
}

func (t *Train) String() string {
	return fmt.Sprintf("Train : %s expected at Station %s at %s leaving at %s",
		t.TrainNumber,
		t.StationUic,
		t.HumanReadableArrivalDate(),
		t.HumanReadableDepartureDate())
}

func (t *Train) arrivesBefore(other *Train) bool {
	return humanReadableDate(t.ScheduledArrival).Before(humanReadableDate(other.ScheduledArrival))
}

func (t *Train) HumanReadableArrivalDate() string {
	return humanReadableDateString(t.ScheduledArrival)
}

func (t *Train) HumanReadableDepartureDate() string {
	return humanReadableDateString(t.ScheduledDeparture)
}

func humanReadableDateString(date string) string {
	return humanReadableDate(date).String()
}

func humanReadableDate(date string) time.Time {
	if strNum := numberSequenceRegex.FindString(date); len(strNum) > 0 {
		if milis, err := strconv.Atoi(strNum); err == nil {
			return time.Unix(int64(milis/1000), 0)
		}
	}
	return time.Time{}
}

type TrainList struct {
	Trains []Train
}

func (tl *TrainList) String() string {
	var buffer bytes.Buffer
	for _, t := range tl.Trains {
		buffer.WriteString(fmt.Sprintf("[ %s ]", t.String()))
	}
	return buffer.String()
}
