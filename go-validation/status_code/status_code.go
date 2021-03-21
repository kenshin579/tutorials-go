package status_code

import (
	"fmt"
	"strconv"
	"strings"
)

//https://gist.github.com/miguelmota/928893046b1fde79fcbe3df52ab332c1
// Code are the status codes that define the stage of the transaction
type Code uint

const (
	// Unprocessed not yet processed
	Unprocessed Code = 1001
	// Pending tx is pending
	Pending Code = 1002
	// Pending tx is awaiting confirmations
	AwaitingConfirmations Code = 1003
	// Failed tx failed
	Failed Code = 5001
	// FeeTooLow tx fee is too low
	FeeTooLow Code = 5002
	// Confirmed tx is confirmed
	Confirmed Code = 12001
)

// Codes are the available codes
var Codes = [6]Code{
	Unprocessed,
	Pending,
	AwaitingConfirmations,
	Failed,
	FeeTooLow,
	Confirmed,
}

// String returns a string representation of the status code
func (o Code) String() string {
	return fmt.Sprintf("string: %d", o)
}

// Message returns the status code message
func (o Code) Message() string {
	switch o {
	case 1001:
		return "UNPROCESSED"

	case 1002:
		return "PENDING"

	case 1003:
		return "AWAITING_CONFIRMATIONS"

	case 5001:
		return "FAILED"

	case 5002:
		return "FEE_TOO_LOW"

	case 12001:
		return "CONFIRMED"

	default:
		return "UNKNOWN"
	}
}

// MarshalJSON marshalls an the status code into json
func (o Code) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, o)), nil
}

// UnmarshalJSON unmarshalls the status code from json
func (o *Code) UnmarshalJSON(data []byte) error {
	if data == nil {
		o = nil
		return nil
	}

	tmp1 := strings.Replace(string(data), `"`, ``, -1)

	tmp, err := strconv.ParseUint(tmp1, 10, 32)
	if err != nil {
		return err
	}
	*o = Code(uint(tmp))
	return nil
}
