package dsb

import (
	"net/http"
)

type APIQuery interface {
	GetFailureChannel() chan error
	GetRequest() *http.Request
	receive([]byte)
}

type query struct {
	failure  chan error
	request  *http.Request
	receiver func([]byte)
}

func NewQuery(failure chan error, request *http.Request, receiver func([]byte)) APIQuery {
	return &query{failure, request, receiver}
}

func (q *query) GetFailureChannel() chan error {
	return q.failure
}

func (q *query) GetRequest() *http.Request {
	return q.request
}

func (q *query) receive(bytes []byte) {
	q.receiver(bytes)
}
