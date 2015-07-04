package dsb

import (
	"encoding/json"
	"fmt"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/utils"
	"io"
	"net/url"
)

func (api *DSBApi) GetTrains() (chan *models.TrainList, chan error) {
	success, failure := make(chan *models.TrainList), make(chan error)
	request, err := api.buildRequest(httpGET, fmt.Sprintf("/Queue()?$filter=(%s) ", url.QueryEscape("TrainType eq 'S-tog'")))
	if err != nil {
		utils.SubmitAsync(err, failure)
	} else {
		q := NewQuery(failure, request, func(body io.Reader) {
			var container map[string][]json.RawMessage
			if err := json.NewDecoder(body).Decode(&container); err != nil {
				failure <- err
			} else {
				success <- convertTrainJSONList(container["d"])
				logger.Println("Completed request")
			}
		})
		api.DoAsync(q)
	}
	return success, failure
}

func convertTrainJSONList(items []json.RawMessage) *models.TrainList {
	result := []models.Train{}
	for _, item := range items {
		m := models.Train{}
		if err := json.Unmarshal(item, &m); err == nil {
			result = append(result, m)
		} else {
			logger.Println(err)
		}
	}
	return &models.TrainList{result}
}
