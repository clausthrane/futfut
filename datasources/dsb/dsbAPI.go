//https://github.com/jordan-wright/gophish
package dsb

import (
	"encoding/json"
	"github.com/clausthrane/futfut/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

type DSBFacade interface {
	GetStations() chan *models.StationList
}

func NewDSBFacade() *DSBApi {
	return &DSBApi{}
}

type DSBApi struct {
}

func (*DSBApi) GetStations() chan *models.StationList {
	replyChan := make(chan *models.StationList)
	go func() {
		client := http.DefaultClient
		req, err := buildRequest()
		resp, err := client.Do(req)
		if err == nil {
			res, _ := doGetStations(resp)
			replyChan <- res
		}
	}()
	return replyChan
}

func buildRequest() (req *http.Request, err error) {
	req, err = http.NewRequest("GET", "http://traindata.dsb.dk/stationdeparture/opendataprotocol.svc/Station()", nil)
	if err == nil {
		req.Header.Add("Accept", "Application/JSON")
	}
	return req, err
}

func doGetStations(resp *http.Response) (*models.StationList, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return unmarshalStations(body), nil
	}
	return nil, err
}

func unmarshalStations(data []byte) *models.StationList {
	stations := []models.Station{}
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

	return &models.StationList{stations}
}

func toStation(data json.RawMessage) (*models.Station, error) {
	station := models.Station{}
	err := json.Unmarshal(data, &station)
	if err == nil {
		return &station, nil
	}
	return nil, err
}
