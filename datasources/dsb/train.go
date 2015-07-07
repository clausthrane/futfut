package dsb

import (
	"encoding/json"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/utils"
	"io"
)

func (api *DSBApi) GetAllTrains() (chan *models.TrainEventList, chan error) {
	return api.getTrainsByQuery("")
}

func (api *DSBApi) GetTrains(key string, value string) (chan *models.TrainEventList, chan error) {
	return api.getTrainsByQuery(filterEQParam(key, value))
}

func (api *DSBApi) getTrainsByQuery(param string) (chan *models.TrainEventList, chan error) {
	success, failure := make(chan *models.TrainEventList), make(chan error)

	request, err := api.buildRequest(httpGET, "/Queue()"+param)

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

func convertTrainJSONList(items []json.RawMessage) *models.TrainEventList {
	result := []models.TrainEvent{}
	for _, item := range items {
		m := models.TrainEvent{}
		if err := json.Unmarshal(item, &m); err == nil {
			result = append(result, m)
		} else {
			logger.Println(err)
		}
	}
	return &models.TrainEventList{result}
}
