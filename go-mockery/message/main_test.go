package main

import (
	"testing"
)

// we need to satisfy our MessageService interface
// which sadly means we have to stub out every method
// defined in that interface
//func (m *smsServiceMock) DummyFunc() {
//	fmt.Println("Dummy")
//}

// TestChargeCustomer is where the magic happens
// here we create our SMSService mock
func TestChargeCustomer(t *testing.T) {
	smsService := &SMSService{}

	// we then define what should be returned from SendChargeNotification
	// when we pass in the value 100 to it. In this case, we want to return
	// true as it was successful in sending a notification
	smsService.On("SendChargeNotification", 100)

	// next we want to define the service we wish to test
	myService := MyService{smsService}
	// and call said method
	myService.ChargeCustomer(100)

	// at the end, we verify that our myService.ChargeCustomer
	// method called our mocked SendChargeNotification method
	smsService.AssertExpectations(t)
}
