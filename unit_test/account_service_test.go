package unit_test_test

import (
	"github.com/stretchr/testify/mock"
)

// MockDMSClient to simulate the dms package's CallAPI function
type MockDMSClient struct {
	mock.Mock
}
