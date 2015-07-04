package models

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
	Line                   string `obj?`
	Direction              string `North, south, ..`
	MinutesToDeparture     string `int`
}

type TrainList struct {
	Trains []Train
}
