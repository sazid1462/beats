// +build !integration

package customendpoint

import (
	"testing"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

type MockTransportClient struct {
	mock.Mock
}

func TestCreateSyslogString(t *testing.T) {

}
