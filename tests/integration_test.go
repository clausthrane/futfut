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

func init() {
	go func() {
		// Start our webapp
		fakeRemote := dsb.NewDSBFacadeWithEndpoint("http://localhost:7777")
		http.ListenAndServe(fmt.Sprintf(":8081"), main.WebAppWithFacade(fakeRemote))
	}()
}

func TestStations(t *testing.T) {
	assert := assert.New(t)

	// well behaved DSB API
	mockserver.HttpServerDSBTestApi(t, 7777)

	if req, err := http.NewRequest("GET", "http://localhost:8081/api/v1/stations", nil); err == nil {
		response, _ := http.DefaultClient.Do(req)
		defer response.Body.Close()
		bytes, _ := ioutil.ReadAll(response.Body)
		jsonString := string(bytes)
		prefix := `{"Count":349,"Stations":[{"Id":"8600020","Name":"Aalborg","CountryCode":"DK"}`
		assert.True(strings.HasPrefix(jsonString, prefix), jsonString[:160])
	} else {
		t.Fail()
	}
}

func TestUnknownResource(t *testing.T) {
	assert := assert.New(t)

	if req, err := http.NewRequest("GET", "http://localhost:8081/api/notexist", nil); err == nil {
		response, _ := http.DefaultClient.Do(req)
		defer response.Body.Close()
		bytes, _ := ioutil.ReadAll(response.Body)
		responseText := string(bytes)
		assert.Equal(responseText, "404 page not found\n", responseText)
	} else {
		t.Fail()
	}
}
