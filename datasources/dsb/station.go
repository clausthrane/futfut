package dsb

import (
	"encoding/json"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/utils"
	"io"
	"sort"
)

func (api *DSBApi) GetStation(stationid string) (chan *models.StationList, chan error) {
	return api.getStationsByQuery(filterEQParam("UIC", stationid))
}

func (api *DSBApi) GetStations() (chan *models.StationList, chan error) {
	return api.getStationsByQuery("")
}

func (api *DSBApi) getStationsByQuery(param string) (chan *models.StationList, chan error) {
	success, failure := make(chan *models.StationList), make(chan error)
	request, err := api.buildRequest(httpGET, "/Station()"+param)
	if err != nil {
		utils.SubmitAsync(err, failure)
	} else {
		q := NewQuery(failure, request, func(body io.Reader) {
			var container map[string][]json.RawMessage
			if err := json.NewDecoder(body).Decode(&container); err != nil {
				failure <- err
			} else {
				success <- convertStationJSONList(container["d"])
			}
		})
		api.DoAsync(q)
	}
	return success, failure
}

func convertStationJSONList(items []json.RawMessage) *models.StationList {
	list := []models.Station{}
	for _, item := range items {
		m := models.Station{}
		if err := json.Unmarshal(item, &m); err == nil {
			list = append(list, m)
		} else {
			logger.Println(err)
		}
	}

	result := &models.StationList{list}
	sort.Sort(result)
	return result
}
