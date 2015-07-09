function FutfutCtrl($scope, $http) {

    startLoading("Warming up")

    $http.get('/api/v1/stations').
        success(function (data) {
            stopLoading()
            $scope.stationCount = data.Count
            $scope.stations = data.Stations;
            $scope.traindata = undefined
        }).error(function (error, status) {
            handleError(error, status)
        });


    $scope.chooseStation = function (stationId) {
        clearAll()
        startLoading("for your station")
        $http.get('/api/v1/departures/from/' + stationId).
            success(function (departures) {
                stopLoading()
                $scope.departures = departures
                $scope.traindata = undefined
            }).error(function (error, status) {
                handleError(error, status)
            });
    }

    $scope.chooseDeparture = function (trainId) {
        removeErrors()
        startLoading("for the train")
        $http.get('/api/v1/trains/' + trainId).
            success(function (trains) {
                stopLoading()
                $scope.trainId = trainId
                $scope.traindata = trains
            }).error(function (error, status) {
                handleError(error, status)
            });
    }

    function startLoading(msg) {
        $scope.loading = msg
    }

    function stopLoading() {
        $scope.loading = undefined
    }

    function removeErrors() {
        $scope.temperror = undefined
        $scope.systemerror = undefined
    }

    function clearAll() {
        stopLoading()
        removeErrors()
        $scope.departures = undefined
        $scope.traindata = undefined
    }

    function handleError(error, status) {
        stopLoading()
        if (status >= 500) {
            $scope.systemerror = error
        } else {
            $scope.temperror = error
        }
    }
}
