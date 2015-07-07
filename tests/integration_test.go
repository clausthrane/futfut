package integrationtest

import (
	"fmt"
	"github.com/clausthrane/futfut"
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/tests/mockserver"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestStations(t *testing.T) {
	assert := assert.New(t)
	StartStack(t)

	if req, err := http.NewRequest("GET", "http://localhost:8081/api/stations", nil); err == nil {
		response, _ := http.DefaultClient.Do(req)
		defer response.Body.Close()
		bytes, _ := ioutil.ReadAll(response.Body)
		prefix := `{"Count":349,"Stations":[{"Id":"7400002","Name":"Göteborg","CountryCode":"S"}`
		assert.True(strings.HasPrefix(string(bytes), prefix))
	} else {
		t.Fail()
	}
}

/*
func TestSpecificStation(t *testing.T) {
	assert := assert.New(t)
	StartStack(t)
	if req, err := http.NewRequest("GET", "http://localhost:8081/api/stations/7400002/details", nil); err == nil {
		response, _ := http.DefaultClient.Do(req)
		defer response.Body.Close()
		bytes, _ := ioutil.ReadAll(response.Body)
		prefix := `{"Id":"7400002","Name":"Göteborg","CountryCode":"S"}`
		assert.True(strings.HasPrefix(string(bytes), prefix))
	} else {
		t.Fail()
	}
}
*/

func StartStack(t *testing.T) {
	mockserver.HttpServerDSBTestApi(t, 7777)
	go func() {
		fakeRemote := dsb.NewDSBFacadeWithEndpoint("http://localhost:7777")
		http.ListenAndServe(fmt.Sprintf(":8081"), main.WebAppWithFacade(fakeRemote))
	}()
}
