function FutfutCtrl($scope, $http) {

    $http.get('/api/v1/stations').
        success(function (data) {
            $scope.stationCount = data.Count
            $scope.stations = data.Stations;
            $scope.traindata = undefined
        });


    $scope.chooseStation = function (stationId) {
        $http.get('/api/v1/departures/from/' + stationId).
            success(function (departures) {
                $scope.departures = departures
                $scope.traindata = undefined
            })
    }

    $scope.chooseDeparture = function (trainId) {
        $http.get('/api/v1/trains/' + trainId).
            success(function (trains) {
                $scope.trainId = trainId
                $scope.traindata = trains
            })
    }
}
