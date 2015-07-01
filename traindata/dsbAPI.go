//https://github.com/jordan-wright/gophish
package traindata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func GetStations() (*StationList, error) {
	client := http.DefaultClient
	req, err := buildRequest()
	resp, err := client.Do(req)
	if err == nil {
		return doGetStations(resp)
	}
	return nil, err
}

func buildRequest() (req *http.Request, err error) {
	req, err = http.NewRequest("GET", "http://traindata.dsb.dk/stationdeparture/opendataprotocol.svc/Station()", nil)
	if err == nil {
		req.Header.Add("Accept", "Application/JSON")
	}
	return req, err
}

func doGetStations(resp *http.Response) (*StationList, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return unmarshalStations(body), nil
	}
	return nil, err
}

func unmarshalStations(data []byte) *StationList {
	stations := []Station{}
	var container map[string][]json.RawMessage
	err := json.Unmarshal(data, &container)
	if err == nil {
		for _, item := range container["d"] {
			if station, err := toStation(item); err == nil {
				stations = append(stations, *station)
			}
		}
	} else {
		logger.Println(err)
	}

	return &StationList{stations}
}

func toStation(data json.RawMessage) (*Station, error) {
	station := Station{}
	err := json.Unmarshal(data, &station)
	if err == nil {
		return &station, nil
	}
	return nil, err
}
