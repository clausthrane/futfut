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

// TrainEvent respresents a train event
//
// E.g. when does a certain train (by TrainNumber) arrive at station Station Uic
type TrainEvent struct {
	ID                     string `UUID`
	StationUic             string `int`
	StationName            string
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

func (t *TrainEvent) String() string {
	return fmt.Sprintf("Train : %s towards %s expected at Station %s at %s leaving at %s",
		t.TrainNumber,
		t.DestinationName,
		t.StationUic,
		t.HumanReadableArrivalDate(),
		t.HumanReadableDepartureDate())
}

func (t *TrainEvent) ArrivesBefore(other *TrainEvent) bool {
	return humanReadableDate(t.ScheduledArrival).Before(humanReadableDate(other.ScheduledArrival))
}

func (t *TrainEvent) ScheduledArrivalDate() time.Time {
	return humanReadableDate(t.ScheduledArrival)
}

func (t *TrainEvent) ScheduledDepartureDate() time.Time {
	return humanReadableDate(t.ScheduledDeparture)
}

func (t *TrainEvent) HumanReadableArrivalDate() string {
	return humanReadableDateString(t.ScheduledArrival)
}

func (t *TrainEvent) HumanReadableDepartureDate() string {
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

// TrainList represent a collection of train events
type TrainEventList struct {
	Events []TrainEvent
}

func (tl *TrainEventList) String() string {
	var buffer bytes.Buffer
	for _, t := range tl.Events {
		buffer.WriteString(fmt.Sprintf("[ %s ]", t.String()))
	}
	return buffer.String()
}

// len, swap, less

func (list *TrainEventList) Len() int {
	return len(list.Events)
}

func (list *TrainEventList) Less(i int, j int) bool {
	itemI := list.Events[i]
	itemJ := list.Events[j]
	return itemI.ScheduledDepartureDate().Before(itemJ.ScheduledDepartureDate())
}

func (list *TrainEventList) Swap(i int, j int) {
	list.Events[i], list.Events[j] = list.Events[j], list.Events[i]
}
