package mockdatabase

import "errors"

type Database interface {
	connect() error
	sendMessage(*string) error
}

func Talk(o Database, message *string) error {
	err := o.connect()

	if err != nil {
		return errors.New("Connection failed")
	}

	err = o.sendMessage(message)
	if err != nil {
		return errors.New("Sending message failed")
	}

	return nil
}
