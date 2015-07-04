package dsb

import (
	"io"
	"net/http"
)

type APIQuery interface {
	GetFailureChannel() chan error
	GetRequest() *http.Request
	receive(io.Reader)
}

type query struct {
	failure  chan error
	request  *http.Request
	receiver func(io.Reader)
}

func NewQuery(failure chan error, request *http.Request, receiver func(io.Reader)) APIQuery {
	return &query{failure, request, receiver}
}

func (q *query) GetFailureChannel() chan error {
	return q.failure
}

func (q *query) GetRequest() *http.Request {
	return q.request
}

func (q *query) receive(body io.Reader) {
	q.receiver(body)
}
