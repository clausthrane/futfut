Futfut (Danish Railway Info - Light)
=========

Demo app for getting and displaying railway data from the Danish railway system (DSB). The app exposes real-time info about depatures and train routes leveraging the API at [DSB Labs](http://www.dsb.dk/dsb-labs/webservice-stationsafgange/) The focus of the app is mostly around the backend, whereas the frontend is fairly minimal.

### Back-end

The back-end for this app is written in [GO] (http://golang.org) (first time experience for me) and has the following structure:

* futfut.go -- Main entry point
* config.json -- Configuration of remote hosts etc.
* /rest -- Implementes the rest layer, mapping resource paths to views
* /view -- Handles marshalling and query parameters
* /models -- Contains domain objects *Station* and *TrainEvent*
* /services -- Implements our "business logic"
* /datasources -- Wrappes the DSB [TrainData API] (http://www.dsb.dk/dsb-labs/webservice-stationsafgange/)
* /config -- Handles reading and distribution of configurations
* /utils -- Implements general utilities e.g. caching layers etc.
* /tests -- Implements integration tests and mock servers and test harnesses

The app exposes a REST API under `/api/v1/` (implemented in `rest/api.go`):

* `GET /api/v1/trains/{id}` Returns the route of a given train, and `404` if not found

* `GET /api/v1/stations` Returns all available stations
* `GET /api/v1/stations/{id}/details` Returns station details, and `404` if not found

* `GET /api/v1/departures/` Returns all departures
* `GET /api/v1/departures/from/{id}` Returns departures from a station, and `404` if not found

The API currently only supports JSON and sets CORS headers for all resources.

### Front-end
The frontend is a simple AngularJS app

* /web - stores the all webcontent, which is served on request to `GET /`

Try it out!
-------------
You can deploy it your self or go have a look at the one I have running.

### Hosted

You should be able to find a [running instance] (http://72.2.119.153:8080/) hosted by
Joyent

### Docker

You can easily build your own docker image of this app as so:

* docker build -t futfut .
* docker run --publish 6060:8080 --name test --rm futfut

An image is also being build automatically and pushed to [Docker Hub](https://registry.hub.docker.com) whenever commits are being merged to the *prod* branch via the Github integration.

* https://registry.hub.docker.com/u/clausthrane/futfut/
