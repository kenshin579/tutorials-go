package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/avast/retry-go"
)

func main() {
	err := retry.Do(
		func() error {
			if err := dummy("frank"); err != nil {
				return err
			}

			return nil
		},
		retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: %s\n", n, err)
		}),
		retry.Attempts(3),
	)

	if err != nil {
		fmt.Printf("err:%v\n", err)
	}

}

func dummy(str string) error {
	fmt.Println(str)
	return errors.New(str)
}
