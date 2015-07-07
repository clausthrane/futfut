package dsb

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestQuery(t *testing.T) {
	assert := assert.New(t)

	bailoutChannel := make(chan error)
	request := &http.Request{}

	called := false

	q := NewQuery(bailoutChannel, request, func(b io.Reader) {
		called = true
	})

	assert.Equal(bailoutChannel, q.GetFailureChannel(), "should match")
	assert.Equal(request, q.GetRequest(), "should match")

	q.receive(nil)
	assert.True(called, "the callback was triggered")

}
