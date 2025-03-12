package mockdatabase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (db *MockDatabase) connect() error {
	args := db.Called()
	return args.Error(0)
}

func (db *MockDatabase) sendMessage(message *string) error {
	args := db.Called(message)
	return args.Error(0)
}

func TestSuccess(t *testing.T) {
	db := new(MockDatabase)
	message := "Hello"

	// Set expectations
	db.On("connect").Return(nil)
	db.On("sendMessage", &message).Return(nil)

	err := Talk(db, &message)

	assert.NoError(t, err)
	db.AssertExpectations(t)
}

func TestErrorOnConnect(t *testing.T) {
	db := new(MockDatabase)

	// Set expectations
	db.On("connect").Return(errors.New("Some error"))

	message := "Hello"
	err := Talk(db, &message)

	assert.NotEqual(t, nil, err, "An error is thrown if connection fails")
	db.AssertExpectations(t)
}

func TestErrorOnMessage(t *testing.T) {
	db := new(MockDatabase)
	message := "Hello"

	// Set expectations
	db.On("connect").Return(nil)
	db.On("sendMessage", &message).Return(errors.New("Some error"))

	err := Talk(db, &message)

	assert.NotEqual(t, nil, err, "An error is thrown if sendMessage fails")
	db.AssertExpectations(t)
}
