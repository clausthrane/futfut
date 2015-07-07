function FutfutCtrl($scope, $http) {

    $http.get('/api/stations').
        success(function (data) {
            $scope.stationCount = data.Count
            $scope.stations = data.Stations;
        });


    $scope.chooseStation = function (stationId) {
        $http.get('/api/departures/from/' + stationId).
            success(function (departures) {
                $scope.departures = departures
                $scope.traindata = undefined
            })
    }

    $scope.chooseDeparture = function (trainId) {
        $http.get('/api/trains/' + trainId).
            success(function (trains) {
                $scope.trainId = trainId
                $scope.traindata = trains
            })
    }
}
