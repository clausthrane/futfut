// This package implements the client for the DSB Traindata API
package dsb

import (
	"fmt"
	"github.com/clausthrane/futfut/config"
	"github.com/clausthrane/futfut/datasources"
	"github.com/clausthrane/futfut/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	httpGET = "GET"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var defaultEndPoint = config.GetString("dsb.endPoint")

type DSBFacade interface {
	GetStations() (chan *models.StationList, chan error)
	GetStation(stationId string) (chan *models.StationList, chan error)
	GetAllTrains() (chan *models.TrainEventList, chan error)
	GetTrains(key string, value string) (chan *models.TrainEventList, chan error)
}

func NewDSBFacade() *DSBApi {
	return &DSBApi{defaultEndPoint}
}

func NewDSBFacadeWithEndpoint(endPoint string) *DSBApi {
	return &DSBApi{endPoint}
}

type DSBApi struct {
	dsbEndpoint string
}

func (api *DSBApi) setEndpoint(endPoint string) {
	api.dsbEndpoint = endPoint
}

func (api *DSBApi) buildRequest(method string, path string) (req *http.Request, err error) {
	logger.Printf("Using DSB endpoint: %s", api.dsbEndpoint)
	req, err = http.NewRequest(method, (*api).dsbEndpoint+path, nil)
	if err == nil {
		req.Header.Add("Accept", "Application/JSON")
	}
	return req, err
}

func (api *DSBApi) DoAsync(q APIQuery) {
	go func() {
		logger.Printf("Performing request: %s %s", q.GetRequest().Method, q.GetRequest().URL)
		response, responsErr := http.DefaultClient.Do(q.GetRequest())
		if responsErr != nil {
			q.GetFailureChannel() <- responsErr
			return
		}
		handleResponse(q, response)
	}()
}

func handleResponse(query APIQuery, resp *http.Response) {
	handleGenericResponse(query.GetFailureChannel(), resp, func(body io.ReadCloser) {
		defer body.Close()
		query.receive(body)
	})
}

type bodyReader func(io.ReadCloser)

// handleGenericResponse takes care of standard error codes
//
// Checks status codes and outputs on 'fail' when receiving a status >= 300
// otherwise the body of the response is applied to the 'concreteResponseHandler'
func handleGenericResponse(fail chan error, resp *http.Response, bodyHandler bodyReader) {
	logger.Printf("Respons was %d", resp.StatusCode)
	switch {
	case resp.StatusCode >= 500:
		fail <- traindata.NewRemoteError(fmt.Sprintf("Remote system encountered an error: %d", resp.StatusCode))
	case 500 > resp.StatusCode && resp.StatusCode >= 400:
		fail <- traindata.NewClientError(fmt.Sprintf("Remote system rejected request: %d", resp.StatusCode))
	case 300 > resp.StatusCode:
		bodyHandler(resp.Body)
	}
}

func filterParam(key string, value string, op string) string {
	filter := fmt.Sprintf("%s %s '%s'", key, op, value)
	// https://github.com/golang/go/issues/4013
	return fmt.Sprintf("?$filter=%s", strings.Replace(url.QueryEscape(filter), "+", "%20", -1))
}

func filterEQParam(key string, value string) string {
	return filterParam(key, value, "eq")
}
