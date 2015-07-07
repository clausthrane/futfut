//This package wraps the train data API from http://www.dsb.dk/dsb-labs/
package traindata

import (
	"fmt"
)

type RemoteError struct {
	msg string
}

func NewRemoteError(msg string) RemoteError {
	return RemoteError{msg}
}

func (err RemoteError) Error() string {
	return fmt.Sprintf("Error caused by remote system: %s - Error was %s", err.msg)
}

type ClientError struct {
	msg string
}

func NewClientError(msg string) ClientError {
	return ClientError{msg}
}

func (err ClientError) Error() string {
	return fmt.Sprintf("Error contacting remote system: %s - Error was %s", err.msg)
}
