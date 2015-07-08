function FutfutCtrl($scope, $http) {

    console.log("ddd")

    $http.get('/api/v1/stations').
        success(function (data) {
            $scope.stationCount = data.Count
            $scope.stations = data.Stations;
            $scope.traindata = undefined
        }).error(function (error, status) {
            handleError(error, status)
        });


    $scope.chooseStation = function (stationId) {
        $scope.error = undefined
        $http.get('/api/v1/departures/from/' + stationId).
            success(function (departures) {
                $scope.departures = departures
                $scope.traindata = undefined
            }).error(function (error) {
                console.log("sss")
                $scope.error = error
            });
    }

    $scope.chooseDeparture = function (trainId) {
        $http.get('/api/v1/trains/' + trainId).
            success(function (trains) {
                $scope.trainId = trainId
                $scope.traindata = trains
            }).error(function (error) {
                console.log("sss")
                $scope.error = error
            });
    }

    function handleError(error, status) {
        if (status >= 500) {
            $scope.systemerror = error
        } else {
            $scope.temperror = error
        }
    }
}
