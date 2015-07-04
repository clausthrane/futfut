package dsb

import (
	"encoding/json"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/utils"
)

func (api *DSBApi) GetStations() (chan *models.StationList, chan error) {
	success, failure := make(chan *models.StationList), make(chan error)
	request, err := api.buildRequest(httpGET, "/Station()")
	if err != nil {
		utils.SubmitAsync(err, failure)
	} else {
		q := NewQuery(failure, request, func(data []byte) {
			var container map[string][]json.RawMessage
			if err := json.Unmarshal(data, &container); err != nil {
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
	result := []models.Station{}
	for _, item := range items {
		m := models.Station{}
		if err := json.Unmarshal(item, &m); err == nil {
			result = append(result, m)
		} else {
			logger.Println(err)
		}
	}
	return &models.StationList{result}
}
