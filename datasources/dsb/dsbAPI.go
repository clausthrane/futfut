//https://github.com/jordan-wright/gophish
package dsb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clausthrane/futfut/config"
	"github.com/clausthrane/futfut/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)
var defaultEndPoint = config.GetString("dsb.endPoint")

func init() {
	logger.Printf("Using DSB endpoint: %s", defaultEndPoint)
}

type DSBFacade interface {
	GetStations() (chan *models.StationList, chan error)
	SetEndpoint(string)
}

func NewDSBFacade() *DSBApi {
	return &DSBApi{defaultEndPoint}
}

type DSBApi struct {
	dsbEndpoint string
}

func (api *DSBApi) SetEndpoint(endPoint string) {
	api.dsbEndpoint = endPoint
}

func (api *DSBApi) GetStations() (chan *models.StationList, chan error) {
	replyChan := make(chan *models.StationList)
	errChan := make(chan error)
	go func() {
		client := http.DefaultClient
		req, err := api.buildRequest()
		resp, err := client.Do(req)
		if err == nil {
			handleGetStationsResponse(replyChan, errChan, resp)
		} else {
			errChan <- err
		}
	}()
	return replyChan, errChan
}

func (api *DSBApi) buildRequest() (req *http.Request, err error) {
	logger.Printf("Preparing request to: %s", api.dsbEndpoint)
	req, err = http.NewRequest("GET", (*api).dsbEndpoint+"/Station()", nil)
	if err == nil {
		req.Header.Add("Accept", "Application/JSON")
	}
	return req, err
}

func handleGetStationsResponse(succ chan *models.StationList, fail chan error, resp *http.Response) {
	logger.Printf("Handling %d response", resp.StatusCode)
	switch {
	case resp.StatusCode >= 500:
		fail <- errors.New(fmt.Sprintf("remote resource is unavailable: %d", resp.StatusCode))
	case 500 > resp.StatusCode && resp.StatusCode >= 400:
		fail <- errors.New(fmt.Sprintf("Internal Server Error: %d", resp.StatusCode))
	case 300 > resp.StatusCode:
		{
			res, marshalError := handleGetStationsBody(resp.Body)
			if marshalError == nil {
				succ <- res
			} else {
				fail <- marshalError
			}
		}
	}
}

func handleGetStationsBody(body io.ReadCloser) (*models.StationList, error) {
	defer body.Close()
	bytes, err := ioutil.ReadAll(body)
	if err == nil {
		return unmarshalStations(bytes)
	}
	return nil, err
}

func unmarshalStations(data []byte) (*models.StationList, error) {
	stations := []models.Station{}
	var container map[string][]json.RawMessage
	err := json.Unmarshal(data, &container)
	if err == nil {
		for _, item := range container["d"] {
			if station, err := toStation(item); err == nil {
				stations = append(stations, *station)
			} else {
				// skipping element; hoping for a partial result
				logger.Println(err)
			}
		}
		return &models.StationList{stations}, nil
	} else {
		return nil, err
	}

}

func toStation(data json.RawMessage) (*models.Station, error) {
	station := models.Station{}
	err := json.Unmarshal(data, &station)
	if err == nil {
		return &station, nil
	}
	return nil, err
}
